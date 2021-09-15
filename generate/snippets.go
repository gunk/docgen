package generate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kenshaw/snaker"
)

var codeTags = map[string]string{
	"sh":  "shell",
	"cpp": "c++",
}

func getGlob(name string) ([]string, error) {
	var res []string
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasPrefix(info.Name(), snaker.CamelToSnake(name)+".") {
				res = append(res, path)
			}

			return nil
		})
	return res, err
}

func addSnippet(name string, langs []string) (string, error) {
	matches, err := getGlob(name)
	if err != nil {
		return "", err
	}
	var res bytes.Buffer
	for _, ext := range langs {
		ext = strings.TrimSpace(ext)
	MATCHES:
		for _, file := range matches {
			// node, err := ci.createNode(v, ext, ignoreExt)
			tag := strings.TrimPrefix(filepath.Ext(file), ".")
			if tag != ext {
				continue MATCHES
			}
			if _, ok := codeTags[tag]; ok {
				tag = codeTags[tag]
			}
			inject, err := ioutil.ReadFile(file)
			if err != nil {
				return "", fmt.Errorf("read code file %s: %w", file, err)
			}
			if tag != "md" {
				// fake source here
				source := "```" + tag + "\n"
				source += string(inject)
				if inject[len(inject)-1] != '\n' {
					source += "\n"
				}
				source += "```"
				inject = []byte(source)
			}
			res.Write(inject)
			res.WriteString("\n\n")
		}
	}
	return res.String(), nil
}
