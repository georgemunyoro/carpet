package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type DirectoryMapEntity struct {
	Filepath    string
	Name        string
	IsDirectory bool
	Contents    map[string]DirectoryMapEntity
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseHeader(t string) interface{} {
	lines := strings.Split(t, "\n")

	inHeader := false
	jStr := ""
	for _, j := range lines {
		if j == "---" {
			inHeader = !inHeader
			if !inHeader {
				break
			}
			continue
		}
		jStr += j + "\n"
	}

	var res interface{}
	err := json.Unmarshal([]byte(jStr), &res)
	check(err)

	return res
}

func extractContent(t string) string {
	lines := strings.Split(t, "\n")
	inHeader := false
	res := ""
	for _, j := range lines {
		if j == "---" {
			inHeader = !inHeader
			continue
		}

		if !inHeader {
			res += j + "\n"
		}
	}
	return res
}

func readFile(filename string) (string, interface{}) {
	dat, err := ioutil.ReadFile(filename)
	check(err)

	if len(strings.TrimSpace(string(dat))) == 0 {
		return "", nil
	}

	content := extractContent(string(dat))
	header := parseHeader(string(dat))

	return content, header
}
