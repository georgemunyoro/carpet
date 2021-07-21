package main

import (
	"testing"
)

var ExamplePage string = `
---
{
	"title": "Example Page",
	"description": "Lorem ipsum dolor sit amet",
	"author": "John Doe",
	"pi": 3.1415926
}
---
<p>Hello, World!</p>`

func TestParseHeader(t *testing.T) {
	got := parseHeader(ExamplePage).(map[string]interface{})
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
