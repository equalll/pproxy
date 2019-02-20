package watch
import "github.com/equalll/mydebug"

/**
*@see https://github.com/gophertown/looper/blob/master/watch.go
 */
import (
	"errors"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
)

func debugMessage(format string, a ...interface{}) {mydebug.INFO()
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
}

func debugError(msg error) {mydebug.INFO()
	fmt.Println(msg.Error())
}

//type FileNotifyEvent fsnotify.FileEvent

type RecursiveWatcher struct {
	Watcher     *fsnotify.Watcher
	Flags       uint32
	OnEventFunc func(event *fsnotify.FileEvent)
	path        string
	Filter      func(path string) bool
}

func NewRecurisveWatcher(path string) (*RecursiveWatcher, error) {mydebug.INFO()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	rw := &RecursiveWatcher{
		Watcher: watcher,
		path:    path,
		Flags:   fsnotify.FSN_ALL,
	}
	rw.Filter = func(path string) bool {
		return true
	}

	return rw, nil
}
func (rw *RecursiveWatcher) Close() error {mydebug.INFO()
	return rw.Watcher.Close()
}

func (rw *RecursiveWatcher) watchFolder(folder string) {mydebug.INFO()
	err := rw.Watcher.WatchFlags(folder, rw.Flags)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}
	fmt.Println("watch:", folder, rw.Flags)
}

func (rw *RecursiveWatcher) AddFolder(path string) {mydebug.INFO()
	folders := rw.walkSubfolders(path)
	if len(folders) == 0 {
		errors.New("No folders to watch.")
		return
	}
	for _, folder := range folders {
		rw.watchFolder(folder)
	}
}

func (rw *RecursiveWatcher) Run(debug bool) {mydebug.INFO()
	rw.AddFolder(rw.path)
	//        go func() {
	for {
		select {
		case event := <-rw.Watcher.Event:
			if !rw.Filter(event.Name) {
				if debug {
					debugMessage("Skip %s", event)
				}
				continue
			}
			// create a file/directory
			if event.IsCreate() {
				fi, err := os.Stat(event.Name)
				if err != nil {
					// eg. stat .subl513.tmp : no such file or directory
					if debug {
						debugError(err)
					}
				} else if fi.IsDir() {
					if debug {
						debugMessage("Detected new directory %s", event.Name)
					}
					rw.watchFolder(event.Name)
				}
			}
			debugMessage("Detected %s", event)
			rw.OnEventFunc(event)

		case err := <-rw.Watcher.Error:
			log.Println("error", err)
		}
	}
	//        }()
}

// returns a slice of subfolders (recursive), including the folder passed in
func (rw *RecursiveWatcher) walkSubfolders(path string) (paths []string) {mydebug.INFO()
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !rw.Filter(newPath) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			hidden := filepath.HasPrefix(name, ".") && name != "." && name != ".."
			if hidden {
				return filepath.SkipDir
			} else {
				paths = append(paths, newPath)
			}
		}
		return nil
	})
	return paths
}
