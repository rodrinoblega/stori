package uses_cases

import (
	"fmt"
	"os"
)

type ProcessFileUseCase struct {
	fileReader        *FileReaderUseCase
	storeTransactions *StoreTransactionsUseCase
	emailSender       *EmailSummaryUseCase
}

func NewProcessFileUseCase(fileReader *FileReaderUseCase, storeTransactions *StoreTransactionsUseCase, emailSender *EmailSummaryUseCase) *ProcessFileUseCase {
	return &ProcessFileUseCase{
		fileReader:        fileReader,
		storeTransactions: storeTransactions,
		emailSender:       emailSender,
	}
}

func (p *ProcessFileUseCase) Execute(filePath string) error {
	transactions, err := p.fileReader.Execute(filePath)
	if err != nil {
		return err
	}

	err = p.storeTransactions.Execute(transactions)
	if err != nil {
		return err
	}

	err = p.emailSender.Execute(transactions)
	if err != nil {
		return err
	}

	return nil
}

func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	return file, nil
}
