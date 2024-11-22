package uses_cases

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ProcessDirectoryExecutor interface {
	Execute(filePath string) error
}

type ProcessDirectoryUseCase struct {
	processFile ProcessFileExecutor
}

func NewProcessDirectoryUseCase(processFile ProcessFileExecutor) *ProcessDirectoryUseCase {
	return &ProcessDirectoryUseCase{
		processFile: processFile,
	}
}

func (p *ProcessDirectoryUseCase) Execute(directoryPath string) error {
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if p.isCSVFile(info) {
			if err := p.processCSVFile(path, info.Name()); err != nil {
				log.Printf("Error processing file %s: %v\n", info.Name(), err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking through directory %s: %w", directoryPath, err)
	}

	return nil
}

func (p *ProcessDirectoryUseCase) isCSVFile(info os.FileInfo) bool {
	return !info.IsDir() && strings.HasSuffix(info.Name(), ".csv")
}

func (p *ProcessDirectoryUseCase) processCSVFile(filePath, fileName string) error {
	log.Printf("Processing file: %s\n", fileName)
	err := p.processFile.Execute(filePath)
	if err != nil {
		return fmt.Errorf("file not processed: %w", err)
	}

	log.Printf("Successfully processed file: %+v\n", fileName)
	return nil
}
