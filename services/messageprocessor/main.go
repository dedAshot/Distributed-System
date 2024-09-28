package main

import (
	"errors"
	"flag"
	"fmt"
	"messageprocessor/consumer"
	"messageprocessor/store"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	DbCredentials      string
	KafkaUrl           string
	KafkaConsumerGroup string
}

func getConfiguration() (config, error) {
	var config config

	DbCredentials := os.Getenv("DB_CREDENTIALS")
	if DbCredentials == "" {
		return config, errors.New("no db credentials")
	}
	KafkaUrl := os.Getenv("BOOTSTRAP_SERVERS")
	if KafkaUrl == "" {
		return config, errors.New("no kafka BOOTSTRAP_SERVERS")
	}
	KafkaConsumerGroup := os.Getenv("KAFKA_CONSUMER_GROUP")
	if KafkaConsumerGroup == "" {
		return config, errors.New("no kafka consumer group")
	}
	config.DbCredentials = DbCredentials
	config.KafkaUrl = KafkaUrl
	config.KafkaConsumerGroup = KafkaConsumerGroup

	return config, nil
}

func installConnections(config config) error {

	if err := store.New(config.DbCredentials); err != nil {
		return err
	}
	fmt.Println("Connected to DB")

	if err := consumer.New(config.KafkaUrl, config.KafkaConsumerGroup); err != nil {
		return err
	}
	fmt.Println("Connected to Kafka")

	return nil
}

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
	fmt.Println("Go service \"message processor\" started")

	config, err := getConfiguration()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if err = installConnections(config); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
