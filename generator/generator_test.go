package generator

import (
	"fmt"
	"strings"
	"testing"
)

var ExamplePage1 string = `
---
{
	"title": "Example Page",
	"description": "Lorem ipsum dolor sit amet",
	"author": "John Doe",
	"pi": 3.1415926
}
---
<p>Hello, World!</p>`

var ExamplePage2 string = `
---
{
	"title": "Example Page",
	"description": "Lorem ipsum dolor sit amet",
	"author": "John Doe",
	"pi": 3.1415926
}
---
<div class="Playground-controls">
    <select class="Playground-selectExample js-playgroundToysEl" aria-label="Code examples">
        <option value="hello.go">Hello, World!</option>
        <option value="life.go">Conway's Game of Life</option>
        <option value="fib.go">Fibonacci Closure</option>
        <option value="peano.go">Peano Integers</option>
        <option value="pi.go">Concurrent pi</option>
        <option value="sieve.go">Concurrent Prime Sieve</option>
        <option value="solitaire.go">Peg Solitaire Solver</option>
        <option value="tree.go">Tree Comparison</option>
    </select>
    <div class="Playground-buttons">
        <button class="Button Button--primary js-playgroundRunEl" title="Run this code [shift-enter]">Run</button>
        <div class="Playground-secondaryButtons">

            <button class="Button js-playgroundShareEl" title="Share this code">Share</button>
            <a class="Button tour" href="https://tour.golang.org/" title="Playground Go from your browser">Tour</a>

        </div>
    </div>
</div>`

var ExamplePage2Content = strings.TrimSpace(`
<div class="Playground-controls">
    <select class="Playground-selectExample js-playgroundToysEl" aria-label="Code examples">
        <option value="hello.go">Hello, World!</option>
        <option value="life.go">Conway's Game of Life</option>
        <option value="fib.go">Fibonacci Closure</option>
        <option value="peano.go">Peano Integers</option>
        <option value="pi.go">Concurrent pi</option>
        <option value="sieve.go">Concurrent Prime Sieve</option>
        <option value="solitaire.go">Peg Solitaire Solver</option>
        <option value="tree.go">Tree Comparison</option>
    </select>
    <div class="Playground-buttons">
        <button class="Button Button--primary js-playgroundRunEl" title="Run this code [shift-enter]">Run</button>
        <div class="Playground-secondaryButtons">

            <button class="Button js-playgroundShareEl" title="Share this code">Share</button>
            <a class="Button tour" href="https://tour.golang.org/" title="Playground Go from your browser">Tour</a>

        </div>
    </div>
</div>`)

var ExamplePage1Content = strings.TrimSpace("<p>Hello, World!</p>")

var ExampleHeader = map[string]interface{}{
	"title":       "Example Page",
	"description": "Lorem ipsum dolor sit amet",
	"author":      "John Doe",
	"pi":          3.1415926,
}

var TestPageInvalidHeader string = `---
{
	invalid: json
}
---`

func TestRender(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/test.html",
		ProjectDir: "../testing/",
	}
	err := TestPage.NewPage()
	if err != nil {
		t.Error(err)
	}

	got, err2 := TestPage.Render()
	if err2 != nil {
		t.Error(err2)
	}
	wanted := strings.TrimSpace(TestPage.Body)
	got = strings.TrimSpace(got)

	if got != wanted {
		t.Errorf("got %s, wanted %s", got, wanted)
	}
}

func TestRenderEmptyString(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/test.html",
		ProjectDir: "../testing/",
	}
	got, err := TestPage.Render()
	if err != nil {
		t.Error(err)
	}

	if got != "" {
		t.Errorf("got '%s', wanted ''", got)
	}
}

func TestLoadHeader(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/test.html",
		ProjectDir: "../testing/",
	}
	err := TestPage.NewPage()
	if err != nil {
		t.Error(err)
	}

	TestPage.loadHeader(ExamplePage1)
	for k, v := range ExampleHeader {
		if TestPage.Header[k] != v {
			t.Errorf("got %v, wanted %v", TestPage.Header[k], v)
		}
	}
}

func TestLoadInvalidHeader(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/bad_header.html",
		ProjectDir: "../testing/",
	}
	err := TestPage.loadHeader(TestPageInvalidHeader)
	wanted := "invalid character 'i' looking for beginning of object key string"

	if err.Error() != wanted {
		t.Error(fmt.Sprintf("got %s, wanted %s", err.Error(), wanted))
	}
}

func TestLoadBody(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/test.html",
		ProjectDir: "../testing/",
	}
	TestPage.NewPage()

	TestPage.loadBody(ExamplePage1)
	got := strings.TrimSpace(TestPage.Body)
	if got != ExamplePage1Content {
		t.Errorf("got %v, wanted %v", got, ExamplePage1Content)
	}

	TestPage.loadBody(ExamplePage2)
	got = strings.TrimSpace(TestPage.Body)
	if got != ExamplePage2Content {
		t.Errorf("got %v, wanted %v", got, ExamplePage2Content)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/non_existent.html",
		ProjectDir: "../testing/",
	}

	err := TestPage.NewPage()
	wanted := fmt.Sprintf("open %s: no such file or directory", TestPage.Filename)

	if err.Error() != wanted {
		t.Error(fmt.Sprintf("got %s, wanted %s", err.Error(), wanted))
	}
}

func TestRenderNonExistentTemplate(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/test.html",
		ProjectDir: "../testing/",
	}

	err := TestPage.NewPage()
	if err != nil {
		t.Error(err)
	}

	TestPage.Header["template"] = "non_existent"
	_, err = TestPage.Render()
	templateFilename := TestPage.ProjectDir + "templates/" + TestPage.Header["template"].(string) + ".html"
	errorMessage := "unable to read template file: " + templateFilename
	if err.Error() != errorMessage {
		t.Errorf("got '%s', wanted '%s'", err.Error(), errorMessage)
	}
}

func TestLoadHeaderLessFile(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/no_header.html",
		ProjectDir: "../testing/",
	}
	err := TestPage.NewPage()
	wanted := fmt.Sprintf("path property not set in page header of file: %s", TestPage.Filename)

	if err.Error() != wanted {
		t.Error(fmt.Sprintf("got '%s', wanted '%s'", err.Error(), wanted))
	}
}

func TestLoadFile(t *testing.T) {
	TestPage := Page{
		Filename:   "../testing/pages/test.html",
		ProjectDir: "../testing/",
	}
	err := TestPage.NewPage()
	if err != nil {
		return
	}

	err = TestPage.loadFile()
	if err != nil {
		t.Error(err)
	}

	gotBody := strings.TrimSpace(TestPage.Body)
	if gotBody != ExamplePage2Content {
		t.Errorf("got %v, wanted %v", gotBody, ExamplePage2Content)
	}

	for k, v := range TestPage.Header {
		if TestPage.Header[k] != v {
			t.Errorf("got %v, wanted %v", TestPage.Header[k], v)
		}
	}

	TestPage = Page{
		Filename:   "../testing/pages/empty_test.html",
		ProjectDir: "../testing/",
	}
	err = TestPage.NewPage()
	if err != nil {
		return
	}

	err = TestPage.loadFile()
	if err != nil {
		t.Error(err)
	}

	gotBody = strings.TrimSpace(TestPage.Body)
	if gotBody != "" {
		t.Errorf("got %v, wanted %v", gotBody, "")
	}

	if TestPage.Header != nil {
		t.Errorf("got %v, wanted %v", TestPage.Header, nil)
	}
}
