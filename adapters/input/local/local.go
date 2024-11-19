package local

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
)

type LocalSource struct {
	Directory string
}

func (l *LocalSource) WatchDirectory() error {
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
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("error watching directory")
			}
			log.Printf("Error watching directory: %v", err)
		}
	}
}

func (l *LocalSource) GetFileContent(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
