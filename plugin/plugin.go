// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
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

	var card CardData
	card.Name = args.ReportName
	card.Reports = []ReportData{}

	for _, suite := range suites {
		card.Reports = append(card.Reports, ReportData{
			Name: suite.Name,
			Tests: TestData{
				Passed: int64(suite.Totals.Passed),
				Failed: int64(suite.Totals.Failed),
				Errored: int64(suite.Totals.Error),
				Skipped: int64(suite.Totals.Skipped),
				Total: int64(suite.Totals.Tests),
			},
			Time: suite.Totals.Duration.String(),
		})
	}

	writeCard(args.Pipeline.Card.Path, "https://rohit-gohri.github.io/drone-junit/cards/v0Card.json", card)

	return nil
}
