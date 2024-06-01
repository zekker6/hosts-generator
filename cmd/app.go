package cmd

import (
	"context"
	"reflect"
	"sort"
	"time"

	"hosts-generator/cmd/file_writer"
	"hosts-generator/cmd/generator"
	"hosts-generator/cmd/parsers"
)

type App struct {
	clients    []parsers.Parser
	writer     *file_writer.Writer
	lineEnding string
	targetIP   string

	syncPeriod time.Duration

	enableWatch bool

	logger func(fmt string, params ...interface{})
}

func NewApp(clients []parsers.Parser, writer *file_writer.Writer, lineEnding string, targetIP string, syncPeriod time.Duration, enableWatch bool, logger func(fmt string, params ...interface{})) *App {
	return &App{clients: clients, writer: writer, lineEnding: lineEnding, targetIP: targetIP, syncPeriod: syncPeriod, enableWatch: enableWatch, logger: logger}
}

func (a *App) Run(ctx context.Context) error {
	var prevHosts []string
	t := time.NewTicker(a.syncPeriod)
	select {
	case <-ctx.Done():
		return nil

	case <-t.C:
		hosts, err := a.GetHosts()
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(prevHosts, hosts) {
			err := a.WriteHosts(hosts)
			if err != nil {
				panic(err)
			}

			prevHosts = hosts

			if a.logger != nil {
				a.logger("updated hosts file, new hosts: %+v", hosts)

			}
		} else {
			if a.logger != nil {
				a.logger("hosts didn't change, skipping")
			}
		}

		if !a.enableWatch {
			break
		}
	}

	return nil
}

func (a *App) Stop() error {
	return a.writer.Clear()
}

func (a *App) WriteHosts(hosts []string) error {
	fileContent := generator.GenerateStrings(hosts, a.targetIP, a.lineEnding)

	return a.writer.WriteToHosts(fileContent)
}

func (a *App) GetHosts() ([]string, error) {
	hosts := make([]string, 0)

	for _, c := range a.clients {
		clientHosts, err := c.Get()
		if err != nil {
			return []string{}, err
		}

		hosts = append(hosts, clientHosts...)
	}

	sort.Strings(hosts)

	return hosts, nil
}
