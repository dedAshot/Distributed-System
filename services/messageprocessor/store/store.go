package store

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var queryTemplates = make(map[string]string)

func init() {
	htmlDirs, _ := os.ReadDir("./store/queries/")

	for _, entry := range htmlDirs {
		name := entry.Name()
		templ, err := getQueryTemplateFromFile("./store/queries/" + name)
		if err != nil {
			fmt.Println("query templ file"+name+"parse err:", err)
			os.Exit(1)
		}
		queryTemplates[name] = templ
	}
}

func getQueryTemplateFromFile(name string) (string, error) {

	query := ""
	fd, err := os.Open(name)
	if err != nil {
		return "", err
	}

	scan := bufio.NewScanner(fd)

	for scan.Scan() {
		query += scan.Text()
	}

	return query, nil
}

func New(dbCredentials string) error {
	
	err := connectToDb(dbCredentials)
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
