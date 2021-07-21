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
		</div>
`
func TestParseHeader(t *testing.T) {
	got  := parseHeader(ExamplePage1).(map[string]interface{})
	want := map[string]interface{}{
		"title": "Example Page",
		"description": "Lorem ipsum dolor sit amet",
		"author": "John Doe",
		"pi": 3.1415926,
	}

	for k, v := range want {
		if got[k] != v {
			t.Errorf("got %v, wanted %v", got[k], v)
		}
	}
}

func TestExtractContent(t *testing.T) {
	got  := strings.TrimSpace(extractContent(ExamplePage1))
	want := strings.TrimSpace("<p>Hello, World!</p>")
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

	got  = strings.TrimSpace(extractContent(ExamplePage2))
	want = strings.TrimSpace(`
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
		</div>
	`)
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
