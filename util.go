package gocookie

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func GetSQLite3Schema(path string) (string, error) {
	cookiesDB, err := sql.Open("sqlite3", "file:"+path+"?mode=ro")
	if err != nil {
		return "", fmt.Errorf("opening database: %v", err)
	}
	defer cookiesDB.Close()

	rows, err := cookiesDB.Query(`SELECT name FROM sqlite_master WHERE type='table'`)
	if err != nil {
		return "", fmt.Errorf("querying tables: %v", err)
	}
	defer rows.Close()

	var schema string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return "", fmt.Errorf("scanning table name: %v", err)
		}

		colRows, err := cookiesDB.Query(fmt.Sprintf("PRAGMA table_info(%q)", tableName))
		if err != nil {
			return "", fmt.Errorf("getting table info for %s: %v", tableName, err)
		}
		defer colRows.Close()

		schema += fmt.Sprintf("Table: %s\n", tableName)
		for colRows.Next() {
			var (
				cid     int
				name    string
				colType string
				notNull bool
				dfltVal sql.NullString
				pk      bool
			)
			if err := colRows.Scan(&cid, &name, &colType, &notNull, &dfltVal, &pk); err != nil {
				return "", fmt.Errorf("scanning column info for %s: %v", tableName, err)
			}
			schema += fmt.Sprintf("  Column: %s (%s)\n", name, colType)
		}

		if err := colRows.Err(); err != nil {
			return "", fmt.Errorf("iterating columns for %s: %v", tableName, err)
		}
		schema += "\n"
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterating tables: %v", err)
	}

	if schema == "" {
		return "", errors.New("no tables found in database")
	}

	return schema, nil
}

func DecryptAESCBC(key, iv, encryptPass []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encryptLen := len(encryptPass)

	dst := make([]byte, encryptLen)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, encryptPass)

	length := len(dst)
	unpad := int(dst[length-1])
	if length <= unpad {
		return nil, errors.New("decrypt failed")
	}
	return dst[:(length - unpad)], nil
}

func ChromiumTimeToUnix(expiresUTC int64) time.Time {
	if expiresUTC == 0 {
		return time.Time{}
	}
	unixSec := (expiresUTC / 1e6) - 11644473600
	return time.Unix(unixSec, 0)
}
