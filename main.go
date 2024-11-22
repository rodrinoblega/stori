package main

import (
	"fmt"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/frameworks/aws_cloud"
	"github.com/rodrinoblega/stori/setup"
	"github.com/rodrinoblega/stori/uses_cases"
	"log"
	"os"
)

const Path = "/path"

func main() {
	configureLogFlags()

	envConf := config.Load(os.Getenv("ENV"))
	dependencies := setup.InitializeDependencies(envConf)

	switch envConf.Env {
	case "prod":
		runProd(dependencies)
	case "local":
		runLocal(dependencies)
	default:
		log.Fatalf("invalid environment: %s", envConf.Env)
	}

	fmt.Println("The process has finished")
}

func runProd(dependencies *setup.AppDependencies) {
	aws_cloud.S3Handler(dependencies)
}

func runLocal(dependencies *setup.AppDependencies) {
	log.Printf("Processing files in the directory")
	if err := dependencies.ProcessDirectory.Execute(Path); err != nil {
		log.Fatalf("Error processing directory files: %v", err)
	}

	log.Printf("Watching path: %s", Path)
	watchDirectory := uses_cases.NewWatchDirectoryUseCase(dependencies.FileWatcher, dependencies.ProcessFile)
	if err := watchDirectory.Execute(); err != nil {
		log.Fatalf("Error watching directory: %v", err)
	}
}

func configureLogFlags() {
	log.SetPrefix("[Stori-App] ")
	log.SetFlags(log.Ldate | log.Ltime)
}
