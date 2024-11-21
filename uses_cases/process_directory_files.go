package uses_cases

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ProcessDirectoryFilesUseCase struct {
	processFile *ProcessFileUseCase
}

func NewProcessDirectoryFilesUseCase(processFile *ProcessFileUseCase) *ProcessDirectoryFilesUseCase {
	return &ProcessDirectoryFilesUseCase{
		processFile: processFile,
	}
}

func (p *ProcessDirectoryFilesUseCase) Execute(directoryPath string) error {
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

func (p *ProcessDirectoryFilesUseCase) isCSVFile(info os.FileInfo) bool {
	return !info.IsDir() && strings.HasSuffix(info.Name(), ".csv")
}

func (p *ProcessDirectoryFilesUseCase) processCSVFile(filePath, fileName string) error {
	log.Printf("Processing file: %s\n", fileName)
	err := p.processFile.Execute(filePath)
	if err != nil {
		return fmt.Errorf("file not processed: %w", err)
	}

	log.Printf("Successfully processed file: %+v\n", fileName)
	return nil
}
