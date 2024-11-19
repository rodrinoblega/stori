package main

import (
	"github.com/rodrinoblega/stori/adapters/email_sender"
	"github.com/rodrinoblega/stori/adapters/repositories"
	"github.com/rodrinoblega/stori/adapters/watchers"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/uses_cases"
	"log"
	"os"
)

func main() {

	envConf := config.Load(os.Getenv("ENV"))

	database := repositories.New(envConf)

	var inputSource uses_cases.Watcher

	switch envConf.Env {
	case "local":
		inputSource = &watchers.LocalSource{Directory: "/path"}
	default:
		log.Fatalf("invalid environment: %s", envConf.Env)
	}

	processFileUseCase := uses_cases.NewProcessFileUseCase(
		uses_cases.NewFileReaderUseCase(),
		uses_cases.NewStoreTransactionsUseCase(database),
		uses_cases.NewEmailSummaryUseCase(email_sender.NewSMTPEmailSender(
			envConf.EmailHost,
			envConf.EmailPort,
			envConf.EmailUsername,
			envConf.EmailPassword),
			database,
		),
	)

	processDirectoryFiles := uses_cases.NewProcessDirectoryFilesUseCase(processFileUseCase)
	err := processDirectoryFiles.Execute("/path")
	if err != nil {
		log.Fatalf("Error processing files in directory: %s, %v", "/path", err)
	}

	watchDirectory := uses_cases.NewWatchDirectoryUseCase(inputSource, processFileUseCase)
	if err := watchDirectory.Execute(); err != nil {
		log.Fatalf("Error executing watch directory: %v", err)
	}

}
