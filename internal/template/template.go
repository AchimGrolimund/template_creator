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

/* The `embed` package directive is used to include the "imagestream.yaml" file located in the "templates" directory
* into the binary at compile time. This file is then accessible as a string variable in the code.
* This is useful for including static assets like configuration, HTML, CSS, images, or in this case, a YAML template.
* The "imagestream.yaml" file is a template used for creating image streams in the application.
 */
//go:embed templates/imagestream.yaml

var imagestreamTemplate string

/* CreateTemplate is a function that creates a new template based on the provided namespace, type, and service.
* It first checks if the namespace ends with "-prod" to set the environment.
* Then, it creates a directory structure based on the namespace, type, environment, and stage.
* It also creates a patches directory and a resources directory within the base directory.
* After that, it parses the imagestream template and creates a new file in the resources directory.
* Finally, it executes the template with the provided service and namespace and writes the output to the file.
 */
func CreateTemplate(namespace, _type, service string) error {
	// Initialize a new logger
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	// Check if namespace ends with -prod and set the environment accordingly
	var env string
	if strings.HasSuffix(namespace, "-prod") {
		env = "prod"
	} else {
		env = "nonprod"
	}

	// Log the set environment and namespace
	log.Info("Setting environment",
		zap.String("environment", env),
		zap.String("namespace", namespace),
	)

	// Create directory structure based on the namespace, type, environment, and stage
	lastIndex := strings.LastIndex(namespace, "-")
	namespaceWithoutSuffix := ""
	stage := ""
	if lastIndex != -1 {
		namespaceWithoutSuffix = namespace[:lastIndex]
		stage = namespace[lastIndex+1:]
	}

	// Log the created directory structure
	log.Info("Creating directory structure",
		zap.String("namespaceWithoutSuffix", namespaceWithoutSuffix),
		zap.String("stage", stage),
	)
	baseDir := fmt.Sprintf("./%s-%s-%s/%s", _type, namespaceWithoutSuffix, env, stage)

	// Create a patches directory within the base directory
	patchesDir := fmt.Sprintf("%s/patches", baseDir)
	err = os.MkdirAll(patchesDir, os.ModePerm)
	if err != nil {
		log.Error("Failed to create patches directory", zap.Error(err))
		return nil
	}
	log.Info("Created patches directory", zap.String("directory", patchesDir))

	// Create a resources directory within the base directory
	resourceesDir := fmt.Sprintf("%s/resources", baseDir)
	err = os.MkdirAll(resourceesDir, os.ModePerm)
	if err != nil {
		log.Error("Failed to create resources directory", zap.Error(err))
		return nil
	}
	log.Info("Created resources directory", zap.String("directory", resourceesDir))

	// Parse the imagestream template
	tmpl, err := template.New("imagestream").Parse(imagestreamTemplate)
	if err != nil {
		log.Error("Failed to parse template", zap.Error(err))
		return nil
	}
	log.Info("Successfully parsed the imagestream template")

	// Create a new file in the resources directory
	imagestreamFile := fmt.Sprintf("%s/resources/imagestream_%s.yaml", baseDir, service)
	f, err := os.Create(imagestreamFile)
	if err != nil {
		log.Error("Failed to create file", zap.Error(err))
		return nil
	}
	log.Info("Created imagestream file", zap.String("file", imagestreamFile))
	defer f.Close()

	// Execute the template with the provided service and namespace and write the output to the file
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
