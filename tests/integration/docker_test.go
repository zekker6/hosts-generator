//go:build integration
// +build integration

package main

import (
	"hosts-generator/cmd/parsers/traefik_v2"
	"testing"
)

func TestLoadingData(t *testing.T) {
	t.Run("parses services names", func(t *testing.T) {
		cl := traefik_v2.NewTraefikV2Client("http://localhost:8888/api")

		res, err := cl.Get()

		if err != nil {
			t.Errorf("failed to load data from traefik: %+v", err)
		}

		if len(res) == 0 {
			t.Errorf("traefik loaded 0 routes, expected to have at least 1, got: %d", len(res))
		}
	})
}
