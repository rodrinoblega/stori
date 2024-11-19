package uses_cases

import (
	"log"
)

type Watcher interface {
	WatchDirectory(processFile *ProcessFileUseCase) error
}

type WatchDirectoryUseCase struct {
	watcher     Watcher
	processFile *ProcessFileUseCase
}

func NewWatchDirectoryUseCase(watcher Watcher, processFile *ProcessFileUseCase) *WatchDirectoryUseCase {
	return &WatchDirectoryUseCase{
		watcher:     watcher,
		processFile: processFile,
	}
}

func (w *WatchDirectoryUseCase) Execute() error {
	if err := w.watcher.WatchDirectory(w.processFile); err != nil {
		return err
	}
	log.Println("Processing completed")
	return nil
}
