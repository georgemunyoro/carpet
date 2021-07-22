package utils

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

func ParseHeader(t string) interface{} {
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
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func ExtractContent(t string) string {
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

func ReadFile(filename string) (string, interface{}) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if len(strings.TrimSpace(string(dat))) == 0 {
		return "", nil
	}

	content := ExtractContent(string(dat))
	header := ParseHeader(string(dat))

	return content, header
}
