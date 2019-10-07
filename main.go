package main

import (
	"astaxie/flatmap"
	"encoding/json"
	"flag"
	"fmt"
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

func main() {
	fmt.Println("helllo")
	// fileName := "/Users/tortlewortle/Library/Application Support/minecraft/saves/Test/stats/5ef828ed-48c2-4a73-95bc-b9fdd4c9e81e.json"
	fileName := flag.String("file", "", "Absolute path to the stats file")
	statsDir := flag.String("dir", "stats", "Directory to store counters in")
	flag.Parse()
	if !fileExists(*fileName) {
		fmt.Println("Stats file does not exist")
		return
	}

	done := make(chan bool)

	if _, err := os.Stat(*statsDir); os.IsNotExist(err) {
		check(os.Mkdir(*statsDir, os.ModePerm))
	}
	exportStats(*fileName, *statsDir)
	go watchForChanges(*fileName)
	go updateCounters(*fileName, *statsDir)

	<-done
}

func watchForChanges(fileName string) {
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

func exportStats(fileName, statsDir string) {
	dat, err := ioutil.ReadFile(fileName)
	check(err)

	statMap := make(map[string]interface{})
	err = json.Unmarshal(dat, &statMap)
	check(err)

	flatMap, err := flatmap.Flatten(statMap)
	check(err)

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
	check(err)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	defer file.Close()

	file.Write([]byte(strconv.Itoa(value)))
}

func updateCounters(fileName, statsDir string) {
	for {
		if modified {
			exportStats(fileName, statsDir)
			modified = false
		}
		time.Sleep(5 * time.Second)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
