package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"statisticsserver/poller"
	"statisticsserver/store"
	"strconv"
)

var (
	defaultPageRows = 100
)

type Message struct {
	Url      string `json:"url"`
	Req_body []byte `json:"req_body,omitempty"`
}

var pageTempl map[string]string = make(map[string]string)

func getHtmlTemplates(path string) {
	htmlDirs, _ := os.ReadDir(path)

	for _, entry := range htmlDirs {
		name := entry.Name()
		templ, err := getHtmlPageFromFile(filepath.Join(path, name))
		if err != nil {
			fmt.Println("pageHtmlTempl file"+name+"parse err:", err)
			os.Exit(1)
		}
		pageTempl[name] = templ
	}
}

func getHtmlPageFromFile(filePath string) (string, error) {

	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)
	var textTemplate string

	for scanner.Scan() {
		textTemplate += scanner.Text() + "\n"
	}

	if textTemplate == "" {
		return "", errors.New("unsuccessful " + filePath + " file scan")
	}

	return textTemplate, nil
}

func homePage(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handle homepage")

	page, ok := pageTempl["homePage.html"]
	if !ok {
		fmt.Fprintf(w, "returning homepage template not found")
		fmt.Println("returning homepage template not found")
	}

	_, err := io.WriteString(w, page)
	if err != nil {
		fmt.Fprintf(w, "returning homepage error")
		fmt.Println("returning homepage error")
	}
}

func apiGetPage(w http.ResponseWriter, req *http.Request) {
	fmt.Println("apiGetPage homepage")

	var stats []*store.StatRow

	url := req.URL.Query()
	rowCount := defaultPageRows
	if tempCount, err := strconv.Atoi(url.Get("rowcount")); err == nil {
		if tempCount < defaultPageRows {
			rowCount = tempCount
		}
	}
	if startId, err := strconv.Atoi(url.Get("startid")); err != nil {
		stats, err = poller.GetPage(-1, rowCount)
		if err != nil {
			fmt.Println("api GetPage err:", err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "server error")
			return
		}
	} else {
		stats, err = poller.GetPage(startId, rowCount)
		if err != nil {
			fmt.Println("api GetPage err:", err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "server error")
			return
		}
	}

	statsJson, err := json.Marshal(stats)
	if err != nil {
		fmt.Println("server marshaling error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "server error")
		return
	}
	fmt.Println(string(statsJson))

	if _, err := w.Write(statsJson); err != nil {

		fmt.Println("server sending err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "server sending error")
		return
	}
}
