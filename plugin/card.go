// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

type TestData struct {
	Failed int64 `json:"failed"`
	Errored int64 `json:"errored"`
	Skipped int64 `json:"skipped"`
	Passed int64 `json:"passed"`
	Total int64 `json:"total"`
}

type ReportData struct {
	Name string `json:"name"`
	Tests TestData `json:"tests"`
	Time string `json:"time"`
}

type CardData struct {
	Name string `json:"name"`
	Reports []ReportData `json:"reports"`
}

func writeCard(path, schema string, card CardData) {
	data, _ := json.Marshal(map[string]interface{}{
		"schema": schema,
		"data":   card,
	})

	logrus.Debug("Card: ", string(data))

	switch {
	case path == "/dev/stdout":
		writeCardTo(os.Stdout, data)
	case path == "/dev/stderr":
		writeCardTo(os.Stderr, data)
	case path != "":
		ioutil.WriteFile(path, data, 0644)
	}
}

func writeCardTo(out io.Writer, data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	io.WriteString(out, "\u001B]1338;")
	io.WriteString(out, encoded)
	io.WriteString(out, "\u001B]0m")
	io.WriteString(out, "\n")
}
