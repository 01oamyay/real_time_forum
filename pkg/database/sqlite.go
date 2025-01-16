package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"rlf/pkg/config"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectSqlte(c *config.Database) (*sql.DB, error) {
	db, err := sql.Open(c.Driver, c.FileName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	if err = makeMigrations(db, c.SchemeDir); err != nil {
		return nil, err
	}

	return db, nil
}

func makeMigrations(db *sql.DB, schemeDir string) error {
	schemes, err := getSchemes(schemeDir)
	if err != nil {
		return err
	}

	for _, scheme := range schemes {
		prep, err := db.Prepare(scheme)
		if err != nil {
			return err
		}
		if _, err = prep.Exec(); err != nil {
			return err
		}
	}
	return nil
}

func getSchemes(schemeDir string) ([]string, error) {
	var schemes []string
	files, err := os.ReadDir(schemeDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fileName := filepath.Join(schemeDir, file.Name())
		data, err := os.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		schemes = append(schemes, string(data))
	}
	return schemes, nil
}
