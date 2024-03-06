package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// InitWatch 文件监控, linux下有点问题
func InitWatch(f func(), files ...string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("NewWatcher fail", "err", err.Error())
		return err
	}
	fileMap := make(map[string]struct{})
	for _, v := range files {
		st, err := os.Lstat(v)
		if err != nil {
			return err
		}
		if st.IsDir() {
			return err
		}
		err = watcher.Add(filepath.Dir(v))
		if err != nil {
			log.Println("Add fail", "err", err.Error())
			return err
		}
		fileMap[v] = struct{}{}
	}
	Info("watch file", "files", fileMap)

	Go(func() {
		defer watcher.Close()
		for {
			select {
			case e, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if _, exist := fileMap[e.Name]; !exist {
					continue
				}
				if e.Has(fsnotify.Write) {
					Info("file change", "name", e.Name, "event", e.Op)
					f()
				}
			case err := <-watcher.Errors:
				fmt.Printf(" %s\n", err.Error())
			}
		}
	})
	return nil
}
