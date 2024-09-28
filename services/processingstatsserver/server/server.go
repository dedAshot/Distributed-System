package server

import (
	"fmt"
	"net/http"
)

var server *http.Server

func setupServer(port, path string) {
	server = &http.Server{}
	server.Addr = ":" + port
	server.Handler = http.DefaultServeMux
	getHtmlTemplates(path)
	http.HandleFunc("GET /{$}", homePage)
	http.HandleFunc("GET /api/getpage/{$}", apiGetPage)
}

func Start(port, path string) error {
	if server != nil {
		return fmt.Errorf("server is already running")
	}
	setupServer(port, path)

	fmt.Println("Server listen: " + port)
	return server.ListenAndServe()
}
