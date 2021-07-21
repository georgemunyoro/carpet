package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

type fileHandler struct {
	root       http.FileSystem
	projectDir string
}

var eventSourceScript string = `
<script>
	let es = new EventSource("/es-subscribe");
	es.addEventListener("reload-event", (e) => {window.location.reload()});
</script>
`

func serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string, redirect bool, projectDir string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var fileContents string
	var fileNotFound bool = true

	if strings.HasSuffix(name, "/") {
		name += "index.html"
	}

	if strings.HasSuffix(name, ".html") {
		filepath := projectDir + name[1:]
		fileContentsByte, err := ioutil.ReadFile(filepath)
		if err == nil {
			fileContents = string(fileContentsByte)
			fileNotFound = false
		}
	}

	if fileNotFound {
		fileContentsByte, err := ioutil.ReadFile(projectDir + "templates/404.html")
		if err == nil {
			fileContents = string(fileContentsByte)
		} else {
			fileContents = "<h1>404</h1><br>Add a <i>404.html</i> into your templates directory to be shown here instead."
		}
	}

	_, err := fmt.Fprint(w, fileContents+eventSourceScript)
	if err != nil {
		return
	}
	return
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	serveFile(w, r, f.root, path.Clean(r.URL.Path), true, f.projectDir)
}

func Serve(filepath string) {
	fs := http.Dir(filepath)
	fh := fileHandler{root: fs}
	fh.projectDir = filepath
	http.Handle("/", &fh)
	log.Fatal(http.ListenAndServe(":8090", nil))
}
