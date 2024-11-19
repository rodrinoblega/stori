package uses_cases

import "log"

type InputSource interface {
	WatchDirectory() error
	GetFileContent(filePath string) ([]byte, error)
}

type ProcessFileUseCase struct {
	input InputSource
}

func NewProcessFileUseCase(input InputSource) *ProcessFileUseCase {
	return &ProcessFileUseCase{input: input}
}

func (u *ProcessFileUseCase) Execute() error {
	if err := u.input.WatchDirectory(); err != nil {
		return err
	}
	log.Println("Processing completed")
	return nil
}
