package main

import (
	"flag"
	logger "log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"traefik-hosts-generator/cmd/api"
	"traefik-hosts-generator/cmd/file_writer"
	"traefik-hosts-generator/cmd/generator"
)

const ApiUrl = "http://localhost:8080/api"
const LocalIP = "127.0.0.1"

var (
	apiUrl    = flag.String("api", ApiUrl, "specify custom traefik API url, example: 'http://127.0.0.1:8080/api'")
	localIP   = flag.String("ip", LocalIP, "specify custom ip to use in hosts file, example: '192.168.33.10'")
	hostsFile = flag.String("file", file_writer.HostsLocation, "specify custom hosts file location, example: '/etc/hosts_custom'")
	watch     = flag.Bool("watch", false, "enable API polling mode: true/false")
	platform  = flag.String("platform", "", "change line-endings style for hosts file, default: '', available: darwin, windows, linux")
	quiet     = flag.Bool("quiet", false, "disable logging")
	period    = flag.Int("freq", 5, "poll every N seconds")
	provider  = flag.String("provider", "docker", "traefik provider to use")
	postfix   = flag.String("postfix", "", "use unique postifix if 2 parallel instances are running")
)

var writer *file_writer.Writer

func main() {
	flag.Parse()

	lineEnding := detectLineEndings()

	writer = file_writer.NewWriter(*hostsFile, lineEnding, *postfix)

	defer func() {
		log("clearing before exit")
		writer.Clear()
	}()

	handleExit()

	var prevHosts []string
	for {
		hosts, err := api.GetHosts(*apiUrl, *provider)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(prevHosts, hosts) {
			err = writeHosts(hosts, lineEnding)
			if err != nil {
				panic(err)
			}

			prevHosts = hosts

			log("updated hosts file, new hosts: %+v", hosts)
		} else {
			log("traefik hosts didn't change, skipping")
		}

		if !*watch {
			break
		}

		time.Sleep(time.Second * time.Duration(*period))
	}
}

func handleExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log("stop signal received")
		writer.Clear()
		os.Exit(0)
	}()

}

func writeHosts(hosts []string, lineEnding string) error {
	fileContent := generator.GenerateStrings(hosts, *localIP, lineEnding)

	return writer.WriteToHosts(fileContent)
}

func log(fmt string, params ...interface{}) {
	if !*quiet {
		if len(params) == 0 {
			logger.Printf(fmt)
		} else {
			logger.Printf(fmt, params...)
		}
	}
}

func detectLineEndings() string {
	lineEnding := generator.LineEnding
	if *platform != "" {
		switch *platform {
		case "linux":
			lineEnding = "\n"
			break
		case "darwin":
			lineEnding = "\n"
			break
		case "windows":
			lineEnding = "\r\n"
			break
		default:
			log("Unknown platform specified: %s, supported: linux, darwin, windows", *platform)
			os.Exit(1)
		}
	}

	return lineEnding
}
