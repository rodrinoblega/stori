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

const Path = "path"

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
	configureLogFlags()

	envConf := config.Load(os.Getenv("ENV"))

	switch envConf.Env {
	case "prod":
		lambda.Start(handler)
	case "local":
		log.Printf("Environment: %s", envConf.Env)
		log.Printf("Initializing the database connection")
		database := repositories.New(envConf)

		log.Printf("Evaluating if there are any files in the directory")
		processDirectoryFiles := initializeProcessDirectoryFiles(database, envConf)
		err := processDirectoryFiles.Execute("/path")
		if err != nil {
			log.Fatalf("Error processing files in directory: %s, %v", "/path", err)
		}

		watcher := watchers.NewWatcherPath(Path)

		log.Printf("Watching path: %s", Path)
		watchDirectory := intializeWatchDirectory(watcher, database, envConf)
		if err := watchDirectory.Execute(); err != nil {
			log.Fatalf("Error executing watch directory: %v", err)
		}
	default:
		log.Fatalf("invalid environment: %s", envConf.Env)
	}

	fmt.Println("The process has finished")
}

func handler(_ context.Context, event S3Event) {
	envConf := config.Load(os.Getenv("ENV"))
	log.Printf("Environment: %s", envConf.Env)

	log.Printf("Initializing the database connection")
	database := repositories.New(envConf)

	log.Printf("Fetching S3 event data")
	result, err := fetchS3Object(event)
	if err != nil {
		log.Fatalf("Failed to fetch object: %v", err)
	}
	defer result.Body.Close()

	log.Printf("Creating a temp file")
	tempFile, err := createTempFile("/tmp/txns.csv", result.Body)
	if err != nil {
		log.Fatalf("error creating temp file: %v", err)
	}

	processFileUseCase := initializeProcessFile(database, envConf)

	log.Printf("Processing the new file")
	if err := processFileUseCase.Execute(tempFile); err != nil {
		log.Fatalf("error processing files in directory: %v", err)
	}

	log.Printf("Process completed")
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

func intializeWatchDirectory(watcher uses_cases.Watcher, database uses_cases.Database, envConf *config.Config) *uses_cases.WatchDirectoryUseCase {
	processFileUseCase := initializeProcessFile(database, envConf)

	return uses_cases.NewWatchDirectoryUseCase(watcher, processFileUseCase)
}

func initializeProcessDirectoryFiles(database uses_cases.Database, envConf *config.Config) *uses_cases.ProcessDirectoryFilesUseCase {
	processFileUseCase := initializeProcessFile(database, envConf)
	return uses_cases.NewProcessDirectoryFilesUseCase(processFileUseCase)
}

func initializeProcessFile(database uses_cases.Database, envConf *config.Config) *uses_cases.ProcessFileUseCase {
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

func configureLogFlags() {
	log.SetPrefix("[Stori-App] ")
	log.SetFlags(log.Ldate | log.Ltime)
}
