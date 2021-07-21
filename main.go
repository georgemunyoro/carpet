package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"weave"
)

var watcher *fsnotify.Watcher

func processFile(filename, projectDir string) error {
	content, header := readFile(filename)
	var ctx map[string]interface{} = header.(map[string]interface{})
	ctx["content"] = content

	template, err := ioutil.ReadFile(projectDir + "templates/" + ctx["template"].(string) + ".html")
	check(err)
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

	if err := filepath.Walk(projectDir+"pages/", watchDirectory); err != nil {
		check(err)
	}
	fmt.Println(projectDir + "pages/")
	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op != fsnotify.Chmod {
					if !strings.HasSuffix(event.Name, "dist") {
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
		fmt.Println(path)
		return watcher.Add(path)
	}
	return nil
}

func main() {
	var projectDirectory = flag.String("p", ".", "Project directory")
	var hotReloadFlag = flag.Bool("w", false, "Watch files")
	flag.Parse()

	projectDir := string(*projectDirectory)
	if projectDir[len(projectDir)-1] != '/' {
		projectDir += "/"
	}

	if !*hotReloadFlag {
		err := compileProject(projectDir)
		check(err)
	} else {
		watchProject(projectDir)
	}
}
