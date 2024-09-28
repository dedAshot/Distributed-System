package main

import (
	"errors"
	"flag"
	"fmt"
	"httphandler/producer"
	"httphandler/server"
	"httphandler/store"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Port                  string
	DbCredentials         string
	KafkaUrl              string
	KafkaProducerSettings string
}

func getConfiguration() (config, error) {
	var config config

	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
	DbCredentials := os.Getenv("DB_CREDENTIALS")
	if DbCredentials == "" {
		return config, errors.New("no db credentials")
	}
	KafkaUrl := os.Getenv("BOOTSTRAP_SERVERS")
	if KafkaUrl == "" {
		return config, errors.New("no kafka url")
	}
	KafkaProducerSettings := os.Getenv("KAFKA_PRODUCER_SETTINGS")
	if KafkaProducerSettings == "" {
		return config, errors.New("no kafka producer settings")
	}

	config.Port = Port
	config.DbCredentials = DbCredentials
	config.KafkaUrl = KafkaUrl
	config.KafkaProducerSettings = KafkaProducerSettings

	return config, nil
}

func installConnections(config config) error {

	if err := store.New(config.DbCredentials); err != nil {
		return err
	}
	fmt.Println("Connected to DB")

	if err := producer.New(config.KafkaUrl, config.KafkaProducerSettings); err != nil {
		return err
	}
	fmt.Println("Connected to Kafka")

	return nil
}

func startServer(port string) error {

	if err := server.Start(port); err != nil {
		return err
	}

	return nil
}

func init() {
	workEnv := flag.Bool("DevEnviroment", false, "Defines is it running on the dev host machine")
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
	fmt.Println("Go service \"httphandler\" started")

	config, err := getConfiguration()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if err = installConnections(config); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if err = startServer(config.Port); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
