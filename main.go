package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rodrinoblega/stori/adapters/email_sender"
	"github.com/rodrinoblega/stori/adapters/repositories"
	"github.com/rodrinoblega/stori/adapters/watchers"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/uses_cases"
	"io"
	"log"
	"os"
)

type S3Event struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func main() {
	envConf := config.Load(os.Getenv("ENV"))
	fmt.Println("Running environment: " + envConf.Env)

	switch envConf.Env {
	case "local":
		database := repositories.New(envConf)

		var inputSource uses_cases.Watcher

		inputSource = &watchers.LocalSource{Directory: "/path"}

		processFileUseCase := initializeProcessFileUseCase(database, envConf)

		processDirectoryFiles := uses_cases.NewProcessDirectoryFilesUseCase(processFileUseCase)
		err := processDirectoryFiles.Execute("/path")
		if err != nil {
			log.Fatalf("Error processing files in directory: %s, %v", "/path", err)
		}

		watchDirectory := uses_cases.NewWatchDirectoryUseCase(inputSource, processFileUseCase)
		if err := watchDirectory.Execute(); err != nil {
			log.Fatalf("Error executing watch directory: %v", err)
		}
	case "prod":
		lambda.Start(handler)
	default:
		log.Fatalf("invalid environment: %s", envConf.Env)
	}

	fmt.Println("The process has finished")
}

func handler(_ context.Context, event S3Event) {
	envConf := config.Load(os.Getenv("ENV"))
	fmt.Println("Running environment: " + envConf.Env)

	database := repositories.New(envConf)

	result, err := fetchS3Object(event)
	if err != nil {
		log.Fatalf("Failed to fetch object: %v", err)
	}
	defer result.Body.Close()

	tempFile, err := createTempFile("/tmp/txns.csv", result.Body)
	if err != nil {
		log.Fatalf("error creating temp file: %v", err)
	}

	processFileUseCase := initializeProcessFileUseCase(database, envConf)

	if err := processFileUseCase.Execute(tempFile); err != nil {
		log.Fatalf("error processing files in directory: %v", err)
	}

	fmt.Println("The lambda process has finished")
}

func extractS3Details(event S3Event) (string, string) {
	return event.Records[0].S3.Bucket.Name, event.Records[0].S3.Object.Key
}

func fetchS3Object(event S3Event) (*s3.GetObjectOutput, error) {
	bucket, key := extractS3Details(event)

	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)
	return s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

func createTempFile(dirPath string, body io.Reader) (string, error) {
	tempFilePath := fmt.Sprintf("%s", dirPath)
	file, err := os.Create(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}

	return tempFilePath, nil
}

func initializeProcessFileUseCase(database uses_cases.Database, envConf *config.Config) *uses_cases.ProcessFileUseCase {
	emailSender := email_sender.NewSMTPEmailSender(
		envConf.EmailHost,
		envConf.EmailPort,
		envConf.EmailUsername,
		envConf.EmailPassword,
	)

	return uses_cases.NewProcessFileUseCase(
		uses_cases.NewFileReaderUseCase(),
		uses_cases.NewStoreTransactionsUseCase(database),
		uses_cases.NewEmailSummaryUseCase(emailSender, database),
	)
}
