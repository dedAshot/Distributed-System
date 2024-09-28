package store

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	_ "github.com/lib/pq"
)

type queryTempl struct {
	TableName string
}

func NewMessage(url string, req_body []byte) *Message {
	return &Message{Url: url, Req_body: req_body}
}

var db *sql.DB

func New(dbCredentials string) error {

	err := connectToDb(dbCredentials)
	for retries := 10; retries > 0 && err != nil; retries-- {
		time.Sleep(time.Second)
		err = connectToDb(dbCredentials)
	}
	if err != nil {
		return err
	}

	exist, err := checkTbleInDb("messages")
	if err != nil {
		return fmt.Errorf("check table in db err:" + err.Error())

	}

	if !exist {
		err = createSchema("messages")

		if err != nil {
			return fmt.Errorf("create table err: %w", err)
		}
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

func checkTbleInDb(tableName string) (bool, error) {
	fmt.Println("Check db existance of table", tableName)

	// templ, err := getQueryTemplateFromFile("checkTablesExistanceTemplate.txt")
	// if err != nil {
	// 	return false, fmt.Errorf("query template parse error: %w", err)
	// }

	// queryData := new(bytes.Buffer)

	// if err = templ.Execute(queryData, queryTempl{TableName: tableName}); err != nil {
	// 	return false, fmt.Errorf("query template execute error: %w", err)
	// }

	fileName := "checkTablesExistanceTemplate.sql"

	file, err := os.Open("./store/" + fileName)

	if err != nil {
		return false, fmt.Errorf("open "+fileName+" error: %w", err)
	}

	scanner := bufio.NewScanner(file)
	var query string

	for scanner.Scan() {
		query += scanner.Text() + "\n"
	}

	var exist bool
	if err := db.QueryRow(query, tableName).Scan(&exist); err != nil {
		var cmpErr string = "pq: column \"" + tableName + "\" does not exist"

		if strings.Compare(err.Error(), cmpErr) == 0 {
			return false, nil
		}
		return false, fmt.Errorf("query execute err: %w", err)
	}

	return exist, nil
}

func getQueryTemplateFromFile(fileName string) (*template.Template, error) {

	file, err := os.Open("./store/" + fileName)

	if err != nil {
		return nil, fmt.Errorf("open template "+fileName+" error: %w", err)
	}

	scanner := bufio.NewScanner(file)
	var textTemplate string

	for scanner.Scan() {
		textTemplate += scanner.Text() + "\n"
	}

	if textTemplate == "" {
		return nil, errors.New("unsuccessful " + fileName + " file scan")
	}

	templ, err := template.New(fileName).Parse(textTemplate)
	if err != nil {
		return nil, fmt.Errorf("template "+fileName+" parse error: %w", err)
	}

	return templ, nil
}

func createSchema(tableName string) error {
	fmt.Println("create table", tableName)

	var query string

	fd, err := os.Open("./store/create_table_" + tableName + ".sql")
	if err != nil {
		return fmt.Errorf("open file "+tableName+".sql error: %w", err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)

	for scanner.Scan() {
		query += scanner.Text() + "\n"
	}

	_, err = db.Exec(query)

	if err != nil {
		return fmt.Errorf("creating "+tableName+" table error: %w", err)
	}

	return nil
}
