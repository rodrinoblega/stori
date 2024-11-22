package uses_cases

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockWatcher struct {
	mock.Mock
}

func (m *MockWatcher) WatchDirectory(processFile *ProcessFileUseCase) error {
	args := m.Called(processFile)
	return args.Error(0)
}

func TestWatchDirectoryUseCase_Execute_Success(t *testing.T) {
	mockWatcher := new(MockWatcher)
	mockProcessFile := new(ProcessFileUseCase)

	useCase := NewWatchDirectoryUseCase(mockWatcher, mockProcessFile)

	mockWatcher.On("WatchDirectory", mockProcessFile).Return(nil)

	err := useCase.Execute()

	assert.NoError(t, err)
	mockWatcher.AssertExpectations(t)
}

func TestWatchDirectoryUseCase_Execute_Failure(t *testing.T) {
	mockWatcher := new(MockWatcher)
	mockProcessFile := new(ProcessFileUseCase)

	useCase := NewWatchDirectoryUseCase(mockWatcher, mockProcessFile)

	mockWatcher.On("WatchDirectory", mockProcessFile).Return(mockError)

	err := useCase.Execute()

	assert.Error(t, err)
	assert.Equal(t, mockError, err)
	mockWatcher.AssertExpectations(t)
}

var mockError = fmt.Errorf("watcher error")
