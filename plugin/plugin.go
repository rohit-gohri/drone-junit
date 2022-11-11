// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshdk/go-junit"
	"github.com/sirupsen/logrus"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	PathsGlob string `envconfig:"PLUGIN_PATHS"`
	ReportName string `envconfig:"PLUGIN_REPORT_NAME"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	globPath := args.PathsGlob
	logrus.Debug("globPath: ", globPath)

	cwd, err := os.Getwd()
	logrus.Debug("CWD: ", cwd)

	if err != nil {
		logrus.Error(err)
	}

	if !strings.HasPrefix(globPath, "/") {
		globPath = cwd + "/" + globPath
	}
	logrus.Debug("Final globPath: ", globPath)

	files, err := filepath.Glob(args.PathsGlob)

	if err != nil {
		logrus.Error(err)
		return errors.New("Invalid glob path")
	}

	if len(files) == 0 {
		return errors.New("No files matched")
	}

	logrus.Debug("Files: ", files)
	
	suites, err := junit.IngestFiles(files)
	
	if err != nil {
		logrus.Error(err)
		return errors.New("Could not parse junit files")
	}

	fmt.Println("***Suites***")
	
	for _, suite := range suites {
		fmt.Println(suite.Name)
		for _, test := range suite.Tests {
			fmt.Printf("  %s\n", test.Name)
			if test.Error != nil {
				fmt.Printf("    %s: %v\n", test.Status, test.Error)
			} else {
				fmt.Printf("    %s\n", test.Status)
			}
		}
	}

	return nil
}
