package store

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var queryTemplates = make(map[string]string)

func chacheQueries(path string) error {
	queryDirs, _ := os.ReadDir(path)

	for _, entry := range queryDirs {
		name := entry.Name()
		templ, err := getQueryTemplateFromFile(filepath.Join(path, name))
		if err != nil {
			fmt.Println("query templ file "+name+" parse err:", err)
			return err
		}
		queryTemplates[name] = templ
	}
	return nil
}

func getQueryTemplateFromFile(name string) (string, error) {

	query := ""
	fd, err := os.Open(name)
	if err != nil {
		return "", err
	}

	scan := bufio.NewScanner(fd)

	for scan.Scan() {
		query += scan.Text() + "\n"
	}

	return query, nil
}

func New(dbCredentials, queryPath string) error {
	if db != nil {
		return nil
	}

	err := chacheQueries(queryPath)
	if err != nil {
		return err
	}

	err = connectToDb(dbCredentials)
	for retries := 10; retries > 0 && err != nil; retries-- {
		time.Sleep(time.Second)
		err = connectToDb(dbCredentials)
	}
	if err != nil {
		return err
	}

	return nil
}

func connectToDb(dbCredentials string) error {
	if db != nil {
		return nil
	}

	var err error
	if db, err = sql.Open("postgres", dbCredentials); err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		db = nil
		return err
	}

	return nil
}
