package uses_cases

import (
	"log"
)

type Watcher interface {
	WatchDirectory(processFile *ProcessFileUseCase) error
}

type WatchFileUseCase struct {
	watcher     Watcher
	processFile *ProcessFileUseCase
}

func NewWatchFileUseCase(watcher Watcher, processFile *ProcessFileUseCase) *WatchFileUseCase {
	return &WatchFileUseCase{
		watcher:     watcher,
		processFile: processFile,
	}
}

func (w *WatchFileUseCase) Execute() error {
	if err := w.watcher.WatchDirectory(w.processFile); err != nil {
		return err
	}
	log.Println("Processing completed")
	return nil
}
