package server

import (
	"fmt"
	"httphandler/producer"
	"net/http"
)

var serverMux *http.ServeMux

func init() {
	serverMux = http.NewServeMux()

	serverMux.HandleFunc("GET /{$}", homePage)
	serverMux.HandleFunc("POST /", handleBlob)

}

func Start(port string) error {
	defer producer.Producer.Close()

	fmt.Println("Server listen: " + port)
	return http.ListenAndServe(":"+port, serverMux)
}
