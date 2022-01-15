package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"

	"github.com/fsnotify/fsnotify"

	"github.com/belljustin/register"
)

func main() {
	args := os.Args[1:]
	pluginPath := args[0]

	// Load all the existing plugins
	files, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		_, err := plugin.Open(filepath.Join(pluginPath, f.Name()))
		if err != nil {
			panic(err)
		}
	}
	displayRates()

	// Watch for new plugins added to the directory
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	watcher.Add(pluginPath)
	defer watcher.Close()

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create != fsnotify.Create {
				continue
			}
			_, err := plugin.Open(event.Name)
			if err != nil {
				panic(err)
			}
			displayRates()
		case err := <-watcher.Errors:
			panic(err)
		}
	}
}

func displayRates() {
	for _, driverName := range register.Registered() {
		forexService := register.Open(driverName)
		rate, err := forexService.GetRate(register.USD, register.CAD)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s:\tUSDCAD\t%d\n", driverName, rate)
	}
}
