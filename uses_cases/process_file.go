package uses_cases

import "github.com/rodrinoblega/stori/entities"

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

func (p *ProcessFileUseCase) Execute(filePath string) (entities.Transactions, error) {
	transactions, err := p.fileReader.Execute(filePath)
	if err != nil {
		return entities.Transactions{}, err
	}

	err = p.storeTransactions.Execute(transactions)
	if err != nil {
		return entities.Transactions{}, err
	}

	err = p.emailSender.Execute(transactions)
	if err != nil {
		return entities.Transactions{}, err
	}

	return transactions, nil
}
