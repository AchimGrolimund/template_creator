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
	"github.com/AchimGrolimund/template_creator/internal/logger"
	"go.uber.org/zap"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/imagestream.yaml
var imagestreamTemplate string

func CreateTemplate(namespace, _type, service string) error {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	// Check if namespace ends with -prod
	var env string
	if strings.HasSuffix(namespace, "-prod") {
		env = "prod"
	} else {
		env = "nonprod"
	}

	// Log with fields
	log.Info("Setting environment",
		zap.String("environment", env),
		zap.String("namespace", namespace),
	)

	// Create directory structure
	lastIndex := strings.LastIndex(namespace, "-")
	namespaceWithoutSuffix := ""
	stage := ""
	if lastIndex != -1 {
		namespaceWithoutSuffix = namespace[:lastIndex]
		stage = namespace[lastIndex+1:]
	}

	// Log with fields
	log.Info("Creating directory structure",
		zap.String("namespaceWithoutSuffix", namespaceWithoutSuffix),
		zap.String("stage", stage),
	)
	baseDir := fmt.Sprintf("./%s-%s-%s/%s", _type, namespaceWithoutSuffix, env, stage)

	patchesDir := fmt.Sprintf("%s/patches", baseDir)
	err = os.MkdirAll(patchesDir, os.ModePerm)
	if err != nil {
		log.Error("Failed to create patches directory", zap.Error(err))
		return nil
	}
	log.Info("Created patches directory", zap.String("directory", patchesDir))

	resourceesDir := fmt.Sprintf("%s/resources", baseDir)
	err = os.MkdirAll(resourceesDir, os.ModePerm)
	if err != nil {
		log.Error("Failed to create resources directory", zap.Error(err))
		return nil
	}
	log.Info("Created resources directory", zap.String("directory", resourceesDir))

	// Render the imagestream template
	tmpl, err := template.New("imagestream").Parse(imagestreamTemplate)
	if err != nil {
		log.Error("Failed to parse template", zap.Error(err))
		return nil
	}
	log.Info("Successfully parsed the imagestream template")

	imagestreamFile := fmt.Sprintf("%s/resources/imagestream_%s.yaml", baseDir, service)
	f, err := os.Create(imagestreamFile)
	if err != nil {
		log.Error("Failed to create file", zap.Error(err))
		return nil
	}
	log.Info("Created imagestream file", zap.String("file", imagestreamFile))
	defer f.Close()

	data := map[string]string{
		"Name":      service,
		"Namespace": namespace,
	}
	err = tmpl.Execute(f, data)
	if err != nil {
		log.Error("Failed to execute template", zap.Error(err))
		return nil
	}
	log.Info("Executed template with data", zap.Any("data", data))

	return nil
}
