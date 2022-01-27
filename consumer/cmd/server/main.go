package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/wednesday-solutions/go-template-consumer/internal/config"
	"github.com/wednesday-solutions/go-template-consumer/pkg/api"
	kafka "github.com/wednesday-solutions/go-template-consumer/pkg/utl/kafkaservice"
)

func main() {
	Setup()
}
func Setup() {
	envName := os.Getenv("ENVIRONMENT_NAME")

	if envName == "" {
		envName = "local"
	}
	kafka.Initiate()
	err := godotenv.Load(fmt.Sprintf(".env.%s", envName))
	if err != nil {
		fmt.Print("error loading .env file")
		checkErr(err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	checkErr(err)
	_, err = api.Start(cfg)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Panic().Msg(err.Error())
	}
}
