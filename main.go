package gomon

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func Init() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	path, err := os.Getwd()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	if err := filepath.Walk(path, watchDir(watcher)); err != nil {
		fmt.Println("ERROR: ", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("EVENT: %#v\n", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("ERROR: ", err)
			}
		}
	}()

	<-done
}

func watchDir(watcher *fsnotify.Watcher) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}

		return nil
	}

}
