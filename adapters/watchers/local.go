package watchers

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rodrinoblega/stori/uses_cases"
	"log"
)

type LocalSource struct {
	Directory string
}

func (l *LocalSource) WatchDirectory(processFile *uses_cases.ProcessFileUseCase) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.New(fmt.Sprintf("error creating file watcher: %v", err))
	}
	defer watcher.Close()

	err = watcher.Add(l.Directory)
	if err != nil {
		return errors.New(fmt.Sprintf("error adding directory to watcher: %v", err))
	}

	log.Printf("Watching directory: %s", l.Directory)

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
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("error watching directory")
			}
			log.Printf("Error watching directory: %v", err)
		}
	}
}
