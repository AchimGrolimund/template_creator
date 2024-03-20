/*
Copyright © 2024 Achim Grolimund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package template

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/imagestream.yaml
var imagestreamTemplate string

func CreateTemplate(namespace, _type, service string) error {
	// Check if namespace ends with -prod
	var env string
	if strings.HasSuffix(namespace, "-prod") {
		env = "prod"
	} else {
		env = "nonprod"
	}

	// Create directory structure
	lastIndex := strings.LastIndex(namespace, "-")
	namespaceWithoutSuffix := ""
	stage := ""
	if lastIndex != -1 {
		namespaceWithoutSuffix = namespace[:lastIndex]
		stage = namespace[lastIndex+1:]
	}
	baseDir := fmt.Sprintf("./%s-%s-%s/%s", _type, namespaceWithoutSuffix, env, stage)
	// fmt.Printf("Creating directory: %s\n", baseDir)

	err := os.MkdirAll(fmt.Sprintf("%s/patches", baseDir), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/resources", baseDir), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Render the imagestream template
	tmpl, err := template.New("imagestream").Parse(imagestreamTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	imagestreamFile := fmt.Sprintf("%s/resources/imagestream_%s.yaml", baseDir, service)
	f, err := os.Create(imagestreamFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	data := map[string]string{
		"Name":      service,
		"Namespace": namespace,
	}
	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
