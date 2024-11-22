package event

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rodrinoblega/stori/setup"
	"github.com/rodrinoblega/stori/src/frameworks/storage"
	"io"
	"log"
	"os"
)

func S3Handler(dependencies *setup.AppDependencies) {
	lambda.Start(func(ctx context.Context, event storage.S3Event) error {
		log.Printf("Received event: %v", event)
		result, err := storage.FetchS3Object(event)
		if err != nil {
			log.Printf("Failed to fetch object: %v", err)
			return fmt.Errorf("failed to fetch object: %w", err)
		}
		defer result.Body.Close()

		tempFilePath, err := createTempFile("/tmp/txns.csv", result.Body)
		if err != nil {
			log.Printf("Failed to create temp file: %v", err)
			return fmt.Errorf("failed to create temp file: %w", err)
		}

		log.Printf("Processing file: %s", tempFilePath)
		if err := dependencies.ProcessFile.Execute(tempFilePath); err != nil {
			log.Printf("Failed to process file: %v", err)
			return fmt.Errorf("failed to process file: %w", err)
		}

		log.Printf("Successfully processed file: %s", tempFilePath)
		return nil
	})
}

func createTempFile(filePath string, body io.Reader) (string, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, body); err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}

	return filePath, nil
}
