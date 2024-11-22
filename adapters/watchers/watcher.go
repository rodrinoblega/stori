package watchers

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rodrinoblega/stori/uses_cases"
	"log"
)

const Path = "/path"

type Watcher struct {
	Directory string
}

func NewWatcherPath(directory string) *Watcher {
	return &Watcher{Directory: directory}
}

func (l *Watcher) WatchDirectory(processFile *uses_cases.ProcessFileUseCase) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.New(fmt.Sprintf("error creating file watcher: %v", err))
	}
	defer watcher.Close()

	err = watcher.Add(l.Directory)
	if err != nil {
		return errors.New(fmt.Sprintf("error adding directory to watcher: %v", err))
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("error detecting new file")
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Printf("New file detected: %s", event.Name)
				err := processFile.Execute(event.Name)
				if err != nil {
					return fmt.Errorf("error creating transaction from csv : %v", err)
				}
				log.Printf("Watching path: %s", Path)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("error watching directory")
			}
			log.Printf("Error watching directory: %v", err)
		}
	}
}
