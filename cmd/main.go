package main

import (
	"User/internal/models"
	"User/internal/server"
	log_c "User/pkg/logger"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log_c.FatalLogger.Fatalf("Error loading .env file: %s", err)
	}

	log_c.InfoLogger.Printf("Successfully loaded environment variables")

	cloudConfig := &models.CloudConfig{
		EndPoint: os.Getenv("AWS_ENDPOINT"),
		Region:   os.Getenv("AWS_REGION"),
		Bucket:   os.Getenv("AWS_S3_BUCKET"),
		Key:      os.Getenv("AWS_KEY"),
		Secret:   os.Getenv("AWS_SECRET"),
	}

	appGRPC, err := server.NewAppGRPC(cloudConfig)
	if err != nil {
		log_c.FatalLogger.Fatalf("Error initializing app: %s", err)
	}

	err = appGRPC.Run(os.Getenv("APP_PORT"))
	if err != nil {
		log_c.FatalLogger.Fatalf("Error running server")
	}
}
