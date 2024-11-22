package uses_cases

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProcessFileUseCase struct {
	mock.Mock
}

func (m *MockProcessFileUseCase) Execute(filePath string) error {
	args := m.Called(filePath)
	return args.Error(0)
}

func TestProcessDirectoryFilesUseCase_Execute_ValidCSVFile(t *testing.T) {
	mockProcessFile := new(MockProcessFileUseCase)
	useCase := NewProcessDirectoryFilesUseCase(mockProcessFile)

	mockProcessFile.On("Execute", "testdata/valid_transactions.csv").Return(nil)

	err := useCase.Execute("testdata/valid_transactions.csv")
	assert.NoError(t, err)

	mockProcessFile.AssertExpectations(t)
}

func TestProcessDirectoryFilesUseCase_Execute_DirectoryError(t *testing.T) {
	mockProcessFile := new(MockProcessFileUseCase)
	useCase := NewProcessDirectoryFilesUseCase(mockProcessFile)

	err := useCase.Execute("invalid/directory/path")
	assert.Error(t, err)
}

func TestProcessDirectoryFilesUseCase_Execute_IgnoreNonCSVFiles(t *testing.T) {
	mockProcessFile := new(MockProcessFileUseCase)
	useCase := NewProcessDirectoryFilesUseCase(mockProcessFile)

	err := useCase.Execute("testdata/other.txt")
	assert.NoError(t, err)

	mockProcessFile.AssertNotCalled(t, "Execute", "testdata/other.txt")
}

func TestProcessDirectoryFilesUseCase_Execute_FileProcessingError(t *testing.T) {
	mockProcessFile := new(MockProcessFileUseCase)
	useCase := NewProcessDirectoryFilesUseCase(mockProcessFile)

	mockProcessFile.On("Execute", "testdata/valid_transactions.csv").Return(errors.New("file processing error"))

	err := useCase.Execute("testdata/valid_transactions.csv")
	assert.Nil(t, err)
}

func TestProcessDirectoryFilesUseCase_Execute_ValidDirectoryAndFile(t *testing.T) {
	mockProcessFile := new(MockProcessFileUseCase)
	useCase := NewProcessDirectoryFilesUseCase(mockProcessFile)

	// Simulate directory with a valid CSV file
	mockProcessFile.On("Execute", "testdata/valid_transactions.csv").Return(nil)

	err := useCase.Execute("testdata/valid_transactions.csv")
	assert.NoError(t, err)
}
