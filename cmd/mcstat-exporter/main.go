package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tortlewortle/mcstat-exporter/internal"
)

func main() {
	fileName := flag.String("file", "", "Absolute path to the stats file")
	statsDir := flag.String("dir", "stats", "Directory to store counters in")
	flag.Parse()

	if !internal.FileExists(*fileName) {
		fmt.Println("Stats file does not exist")
		return
	}

	done := make(chan bool)

	if _, err := os.Stat(*statsDir); os.IsNotExist(err) {
		internal.Check(os.Mkdir(*statsDir, os.ModePerm))
	}
	internal.ExportStats(*fileName, *statsDir)
	go internal.WatchForChanges(*fileName)
	go internal.UpdateCounters(*fileName, *statsDir)

	<-done
}
