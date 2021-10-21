package cmd

import (
	"hosts-generator/cmd/file_writer"
	"hosts-generator/cmd/generator"
	"hosts-generator/cmd/parsers"
	"reflect"
	"time"
)

type App struct {
	clients    []parsers.Parser
	writer     *file_writer.Writer
	lineEnding string
	targetIP   string

	syncPeriod int

	enableWatch bool

	logger func(fmt string, params ...interface{})
}

func NewApp(clients []parsers.Parser, writer *file_writer.Writer, lineEnding string, targetIP string, syncPeriod int, enableWatch bool, logger func(fmt string, params ...interface{})) *App {
	return &App{clients: clients, writer: writer, lineEnding: lineEnding, targetIP: targetIP, syncPeriod: syncPeriod, enableWatch: enableWatch, logger: logger}
}

func (a *App) Run() error {
	var prevHosts []string
	for {
		hosts, err := a.GetHosts()
		if err != nil {
			return nil
		}

		if !reflect.DeepEqual(prevHosts, hosts) {
			err := a.WriteHosts(hosts)
			if err != nil {
				panic(err)
			}

			prevHosts = hosts

			a.logger("updated hosts file, new hosts: %+v", hosts)
		} else {
			a.logger("hosts didn't change, skipping")
		}

		if !a.enableWatch {
			break
		}

		time.Sleep(time.Second * time.Duration(a.syncPeriod))
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
	return hosts, nil
}
