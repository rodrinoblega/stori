package main

import (
	"github.com/rodrinoblega/stori/adapters/watchers"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/frameworks/db"
	"github.com/rodrinoblega/stori/uses_cases"
	"log"
	"os"
)

func main() {

	envConf := config.Load(os.Getenv("ENV"))

	database := db.New(envConf)

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
	)

	useCase := uses_cases.NewWatchFileUseCase(inputSource, processFileUseCase)
	if err := useCase.Execute(); err != nil {
		log.Fatalf("Error executing use case: %v", err)
	}

}
