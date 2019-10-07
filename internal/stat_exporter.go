package internal

import (
	"astaxie/flatmap"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var modified = false

// WatchForChanges watches for stat file changes.
func WatchForChanges(fileName string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(fileName)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
				modified = true
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

// ExportStats exports the stats to .txt files
func ExportStats(fileName, statsDir string) {
	dat, err := ioutil.ReadFile(fileName)
	Check(err)

	statMap := make(map[string]interface{})
	err = json.Unmarshal(dat, &statMap)
	Check(err)

	flatMap, err := flatmap.Flatten(statMap)
	Check(err)
	writeFiles(flatMap, statsDir)
}

func writeFiles(flatMap flatmap.FlatMap, statsDir string) {
	for key, value := range flatMap {
		num, err := strconv.ParseFloat(value, 0)
		if err != nil {
			num = 0
		}
		fileName := strings.Replace(strings.Replace(key+".txt", ":", ".", -1), "/", ".", -1)
		writeFile(statsDir, strings.Replace(fileName, "/", ".", -1), int(num))
	}
}

func writeFile(statsDir, name string, value int) {
	path, err := filepath.Abs(filepath.Join(statsDir, name))
	Check(err)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	Check(err)
	defer file.Close()

	file.Write([]byte(strconv.Itoa(value)))
}

// UpdateCounters updates the counter files
func UpdateCounters(fileName, statsDir string) {
	for {
		if modified {
			ExportStats(fileName, statsDir)
			modified = false
		}
		time.Sleep(5 * time.Second)
	}
}
