package server

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"httphandler/producer"
	"httphandler/store"
	"io"
	"net/http"
	"os"
)

var BLOB_MAX_LENGTH = 10 * 1024 * 1024 //10 MiB

var pageTempl map[string]*template.Template = make(map[string]*template.Template)

func getHtmlTemplateFromFile(filePath string) (*template.Template, error) {

	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	var textTemplate string

	for scanner.Scan() {
		textTemplate += scanner.Text() + "\n"
	}

	if textTemplate == "" {
		return nil, errors.New("unsuccessful" + filePath + "file scan")
	}

	templ, err := template.New(filePath).Parse(textTemplate)
	if err != nil {
		return nil, err
	}

	return templ, nil
}

func init() {
	htmlDirs, _ := os.ReadDir("./server/htmlpagestempl")

	for _, entry := range htmlDirs {
		name := entry.Name()
		templ, err := getHtmlTemplateFromFile("./server/htmlpagestempl/" + name)
		if err != nil {
			fmt.Println("pageHtmlTempl file"+name+"parse err:", err)
			os.Exit(1)
		}
		pageTempl[name] = templ
	}

}

func homePage(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handle homepage")

	fmt.Println(pageTempl)

	templ, ok := pageTempl["homePage.html"]
	if !ok {
		fmt.Fprintf(w, "returning homepage template not found")
		fmt.Println("returning homepage template not found")
	}

	err := templ.Execute(w, nil)
	if err != nil {
		fmt.Fprintf(w, "returning homepage error")
		fmt.Println("returning homepage error")
	}
}

func handleBlob(w http.ResponseWriter, req *http.Request) {

	fmt.Println("handle Blob")

	if req.ContentLength > int64(BLOB_MAX_LENGTH) {
		fmt.Fprintf(w, "Too lot data (>%d Bytes)", BLOB_MAX_LENGTH)
		fmt.Printf("Too lot data (>%d Bytes)\n", BLOB_MAX_LENGTH)
		return
	}

	url := req.URL.Path

	bodyReader := bufio.NewReader(req.Body)
	req_body := make([]byte, 0, 2048)
	buf := make([]byte, 2048)

	if req.ContentLength != 0 {
		n, err := bodyReader.Read(buf)
		for ; err == nil && n != 0; n, err = bodyReader.Read(buf) { // _, err = bodyReader.Read(buf)
			req_body = append(req_body, buf[:n]...)
		}
		if err != io.EOF {
			fmt.Fprintf(w, "An error occured during message processing")
			fmt.Println("An error occured during message processing")
			return
		}
	} else {
		fmt.Fprintf(w, "Empty res body")
		fmt.Println("Empty res body")
		return
	}
	var dbSaveStatus, kafkaSaveStatus string
	msg := store.NewMessage(url, req_body)

	if err := store.MessageRepository.SaveMessage(msg); err != nil {
		fmt.Println(err)
		dbSaveStatus = "db saving error: " + err.Error()
	} else {
		dbSaveStatus = "successfully saved message in db"
	}

	if err := producer.SendMsg(msg); err != nil {
		fmt.Println(err)
		kafkaSaveStatus = "kafka sending error: " + err.Error()
	} else {
		kafkaSaveStatus = "successfully sended message to kafka"
	}

	fmt.Fprintf(w, dbSaveStatus+"\n"+kafkaSaveStatus)
}
