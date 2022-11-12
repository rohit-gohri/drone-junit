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
	card.name = args.ReportName
	card.reports = []ReportData{}

	for _, suite := range suites {
		card.reports = append(card.reports, ReportData{
			name: suite.Name,
			tests: TestData{
				passed: int64(suite.Totals.Passed),
				failed: int64(suite.Totals.Failed),
				errored: int64(suite.Totals.Error),
				skipped: int64(suite.Totals.Skipped),
				total: int64(suite.Totals.Tests),
			},
			time: suite.Totals.Duration.String(),
		})
	}

	writeCard(args.Pipeline.Card.Path, "https://rohit-gohri.github.io/drone-junit/cards/v0Card.json", card)

	return nil
}
