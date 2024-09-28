package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"statisticsserver/poller"
	"statisticsserver/server"
	"statisticsserver/store"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	port   string
	dbcred string
}

var (
	defaultPort = "8080"
	cfg         = &config{}
)

func init() {
	workEnv := flag.Bool("dev", false, "Defines is it running on the dev host machine")
	flag.Parse()

	if *workEnv {
		fmt.Println("dev enviroment")
		if err := godotenv.Load(); err != nil {
			fmt.Println("Err, during get env virables: ", err)
			os.Exit(-1)
		}
	}

}

func main() {
	fmt.Println("http statistic server started")

	getConfig()
	installConnections()
	poller.StartPolling(time.Second)
	templPath := filepath.Join("server", "htmlpagestempl")
	server.Start(cfg.port, templPath)
}

func getConfig() {

	port := os.Getenv("PORT")
	if strings.EqualFold(port, "") {
		port = defaultPort
	}

	dbcred := os.Getenv("DB_CREDENTIALS")
	if strings.EqualFold(dbcred, "") {
		fmt.Println("Empty db creds")
		os.Exit(0)
	}

	cfg.port = port
	cfg.dbcred = dbcred
	fmt.Println(cfg)
}

func installConnections() {
	//store.
	queriesTemplPath := filepath.Join("store", "queries")
	if err := store.New(cfg.dbcred, queriesTemplPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
