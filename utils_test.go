package main

import (
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

func TestParseHeader(t *testing.T) {
	got := parseHeader(ExamplePage1).(map[string]interface{})
	for k, v := range ExampleHeader {
		if got[k] != v {
			t.Errorf("got %v, wanted %v", got[k], v)
		}
	}
}

func TestExtractContent(t *testing.T) {
	got := strings.TrimSpace(extractContent(ExamplePage1))
	if got != ExamplePage1Content {
		t.Errorf("got %v, wanted %v", got, ExamplePage1Content)
	}

	got = strings.TrimSpace(extractContent(ExamplePage2))
	if got != ExamplePage2Content {
		t.Errorf("got %v, wanted %v", got, ExamplePage2Content)
	}
}

func TestReadFile(t *testing.T) {
	gotContent, gotHeader := readFile("./testing/test.html")
	gotContent = strings.TrimSpace(gotContent)
	if gotContent != ExamplePage2Content {
		t.Errorf("got %v, wanted %v", gotContent, ExamplePage2Content)
	}

	for k, v := range ExampleHeader {
		if gotHeader.(map[string]interface{})[k] != v {
			t.Errorf("got %v, wanted %v", gotHeader.(map[string]interface{})[k], v)
		}
	}

	gotContent, gotHeader = readFile("./testing/empty_test.html")
	if gotContent != "" {
		t.Errorf("got %v, wanted %v", gotContent, "")
	}

	if gotHeader != nil {
		t.Errorf("got %v, wanted %v", gotHeader, nil)
	}
}
