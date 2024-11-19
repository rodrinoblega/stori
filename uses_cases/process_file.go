package uses_cases

import "github.com/rodrinoblega/stori/entities"

type ProcessFileUseCase struct {
	fileReader *FileReaderUseCase
}

func NewProcessFileUseCase(fileReader *FileReaderUseCase) *ProcessFileUseCase {
	return &ProcessFileUseCase{
		fileReader: fileReader,
	}
}

func (p *ProcessFileUseCase) Execute(filePath string) (entities.Transactions, error) {
	return p.fileReader.Execute(filePath)
}
