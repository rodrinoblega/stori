package main

import (
	"github.com/rodrinoblega/stori/adapters/input/local"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/uses_cases"
	"log"
	"os"
)

func main() {

	envConf := config.Load(os.Getenv("ENV"))

	var inputSource uses_cases.InputSource

	switch envConf.Env {
	case "local":
		inputSource = &local.LocalSource{Directory: "/path"}
	default:
		log.Fatalf("invalid environment: %s", envConf.Env)
	}

	useCase := uses_cases.NewProcessFileUseCase(inputSource)
	if err := useCase.Execute(); err != nil {
		log.Fatalf("Error executing use case: %v", err)
	}

}
