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
package main

import (
	"github.com/AchimGrolimund/template_creator/cmd"
	"github.com/AchimGrolimund/template_creator/internal/logger"
	"go.uber.org/zap"
	"runtime"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	log.Info("Starting the application")
	// Überprüft, ob das aktuelle Log-Level auf Debug oder niedriger eingestellt ist.
	// Dies ist nützlich, wenn die Erstellung der Debug-Nachricht selbst rechenintensiv ist oder Nebeneffekte hat.
	// In solchen Fällen möchten Sie möglicherweise vermeiden, die Nachricht zu erstellen, wenn sie nicht geloggt wird.
	if log.Core().Enabled(zap.DebugLevel) {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		log.Debug("Debug logging enabled",
			zap.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)),
			zap.Int("NumGoroutine", runtime.NumGoroutine()),
			zap.Int("NumCPU", runtime.NumCPU()),
			zap.Uint64("HeapAlloc", mem.HeapAlloc),
		)
	}
	cmd.Execute()
}
