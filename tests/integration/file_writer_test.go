//go:build integration
// +build integration

package main

import (
	"context"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"hosts-generator/cmd"
	"hosts-generator/cmd/file_writer"
	"hosts-generator/cmd/parsers"
	"hosts-generator/cmd/parsers/traefik_v2"
)

const (
	testHostsFileLocation = "/tmp"
)

func getTempFileName() string {
	f, _ := ioutil.TempFile(testHostsFileLocation, "hosts-generator")

	return f.Name()
}

func TestWritingToFile(t *testing.T) {
	cl := traefik_v2.NewTraefikV2Client("http://localhost:8888/api")

	lineEnding := "\n"

	tmpFile := getTempFileName()
	adapter := file_writer.NewFileHostsAdapter(tmpFile)

	writer := file_writer.NewWriter(&adapter, lineEnding, "")

	app := cmd.NewApp([]parsers.Parser{cl}, writer, lineEnding, "127.0.0.1", 1, false, log.Default().Printf)

	t.Run("runs application", func(t *testing.T) {

		err := app.Run(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %+v", err)
		}
	})

	t.Run("writes are idempotent", func(t *testing.T) {
		for i := 0; i <= 5; i++ {
			initialContent, _ := ioutil.ReadFile(tmpFile)

			err := app.Run(context.Background())
			if err != nil {
				t.Errorf("unexpected error: %+v", err)
			}

			afterRunContent, _ := ioutil.ReadFile(tmpFile)

			if !reflect.DeepEqual(initialContent, afterRunContent) {
				t.Errorf("expected: %s, got: %s", initialContent, afterRunContent)
			}
		}

	})

	t.Run("loads hosts", func(t *testing.T) {
		h, _ := app.GetHosts()
		if len(h) == 0 {
			t.Errorf("expected to find hosts")
		}
	})
}
