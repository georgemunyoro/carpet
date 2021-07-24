package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"weave"
)

type Page struct {
	Filename       string
	OutputFilepath string
	Body           string
	Header         map[string]interface{}
	FileContents   string
	ProjectDir     string
}

func (p *Page) NewPage() error {
	err := p.loadFile()
	if err != nil {
		return err
	}
	if p.Header["path"] == nil {
		return errors.New(fmt.Sprintf(
			"path property not set in page header of file: %s",
			p.Filename,
		))
	}
	p.OutputFilepath = p.ProjectDir + "dist/" + p.Header["path"].(string)
	return nil
}

func (p *Page) Render() (string, error) {
	if strings.TrimSpace(p.Body) == "" {
		return "", nil
	}

	var ctx map[string]interface{} = p.Header
	ctx["content"] = p.Body

	templateFilename := p.ProjectDir + "templates/" + ctx["template"].(string) + ".html"
	template, err := ioutil.ReadFile(templateFilename)
	if err != nil {
		return "", errors.New(fmt.Sprintf(
			"unable to read template file: %s",
			templateFilename,
		))
	}

	if len(template) == 0 {
		return "", nil
	}

	return weave.Render(string(template), ctx), err
}

func (p *Page) WriteToDisk() error {
	//x := strings.Split(p.Filename, "/pages/")[1]
	//y := strings.Split(x, "/")
	//z := y[:len(y)-1]
	//
	//if len(z) > 0 {
	//	err := os.MkdirAll(p.ProjectDir+"dist/"+strings.Join(z, "/"), os.ModePerm)
	//	if err != nil {
	//		return err
	//	}
	//	z = append(z, ctx["path"].(string))
	//	ctx["path"] = strings.Join(z, "/")
	//}

	r, err := p.Render()
	if p.OutputFilepath == "" {
		return errors.New("output filepath not yet generated")
	}
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p.OutputFilepath, []byte(r), os.ModePerm)
}

func (p *Page) loadHeader(s string) error {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	lines := strings.Split(s, "\n")

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

	err := json.Unmarshal([]byte(jStr), &p.Header)
	if err != nil {
		return err
	}

	return nil
}

func (p *Page) loadBody(s string) {
	if strings.TrimSpace(s) == "" {
		return
	}
	lines := strings.Split(s, "\n")
	inHeader := false
	p.Body = ""
	for _, j := range lines {
		if j == "---" {
			inHeader = !inHeader
			continue
		}

		if !inHeader {
			p.Body += j + "\n"
		}
	}
}

func (p *Page) loadFile() error {
	dat, err := ioutil.ReadFile(p.Filename)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(string(dat))) == 0 {
		p.Body = ""
		p.Header = nil
	}

	p.loadBody(string(dat))
	err = p.loadHeader(string(dat))
	if err != nil {
		return err
	}

	return nil
}
