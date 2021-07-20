package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"weave"
)

func processFile(filename, projectDir string) error {
	content, header := readFile(filename)
	var ctx map[string]interface{} = header.(map[string]interface{})
	ctx["content"] = content

	template, err := ioutil.ReadFile(projectDir + "/templates/" + ctx["template"].(string) + ".html")
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

func compileProject(projectDir string) (error) {
	var files []string
	err := filepath.Walk(projectDir + "/pages", func(path string, info os.FileInfo, e error) error {
		check(e)
		if strings.HasSuffix(path, "html") {
			files = append(files, path)
		}
		return nil
	})
	check(err)

	for _, file := range files {
		err := processFile(file, projectDir)
		check(err)
	}
	return nil
}

func main() {
	const testDir = "/home/georgemunyoro/workspace/georgemunyoro.me/"
	err := compileProject(testDir)
	check(err)
}
