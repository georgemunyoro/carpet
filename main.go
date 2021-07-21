package main

import (
	"carpet/server"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/antage/eventsource.v1"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"weave"
)

var watcher *fsnotify.Watcher

func processFile(filename, projectDir string) error {
	content, header := readFile(filename)
	if len(strings.TrimSpace(content)) == 0 {
		return nil
	}
	var ctx map[string]interface{} = header.(map[string]interface{})
	ctx["content"] = content

	template, err := ioutil.ReadFile(projectDir + "templates/" + ctx["template"].(string) + ".html")
	check(err)

	if len(template) == 0 {
		return nil
	}

	render := weave.Render(string(template), ctx)

	x := strings.Split(filename, "/pages/")[1]
	y := strings.Split(x, "/")
	z := y[:len(y)-1]

	if len(z) > 0 {
		err := os.MkdirAll(projectDir+"dist/"+strings.Join(z, "/"), os.ModePerm)
		if err != nil {
		}
		z = append(z, ctx["path"].(string))
		ctx["path"] = strings.Join(z, "/")
	}

	newFilename := projectDir + "dist/" + ctx["path"].(string)
	err = ioutil.WriteFile(newFilename, []byte(render), os.ModePerm)
	return err
}

func compileProject(projectDir string) error {
	var files []string
	err := filepath.Walk(projectDir+"pages", func(path string, info os.FileInfo, e error) error {
		check(e)
		if strings.HasSuffix(path, "html") {
			files = append(files, path)
		}
		return nil
	})
	check(err)

	err = os.RemoveAll(projectDir + "dist")
	check(err)
	err = os.Mkdir(projectDir+"dist", os.ModePerm)
	check(err)

	for _, file := range files {
		if strings.HasSuffix(file, "html") {
			err := processFile(file, projectDir)
			check(err)
		}
	}
	return nil
}

func watchProject(projectDir string) {
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	es := eventsource.New(nil, nil)
	defer es.Close()
	http.Handle("/es-subscribe", es)

	if err := filepath.Walk(projectDir+"pages/", watchDirectory); err != nil {
		check(err)
	}
	if err := filepath.Walk(projectDir+"templates/", watchDirectory); err != nil {
		check(err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op != fsnotify.Chmod {
					if !strings.HasSuffix(event.Name, "dist") && !strings.HasSuffix(event.Name, "templates") {
						es.SendEventMessage("reload", "reload-event", time.Time{}.String())
						fmt.Println("filesystem change detected, reloading...")
						err := compileProject(projectDir)
						check(err)
					}
				}
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	if err := watcher.Add(projectDir); err != nil {
		fmt.Println("ERROR", err)
	}
	<-done
}

func watchDirectory(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() && !strings.HasSuffix(path, ".html") {
		return watcher.Add(path)
	}
	return nil
}

func main() {
	var projectDirectory = flag.String("p", ".", "Project directory")
	var hotReloadFlag = flag.Bool("w", false, "Watch files")
	var serveFlag = flag.Bool("s", false, "Serve output directory")
	flag.Parse()

	projectDir := string(*projectDirectory)
	if projectDir[len(projectDir)-1] != '/' {
		projectDir += "/"
	}

	if *serveFlag {
		go server.Serve(projectDir + "dist/")
	}

	if !*hotReloadFlag {
		err := compileProject(projectDir)
		check(err)
	} else {
		watchProject(projectDir)
	}
}
