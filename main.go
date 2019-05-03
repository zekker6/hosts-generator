package main

import (
	"flag"
	"gitlab.r-styleserv.com/r-style/devops/traefik-hosts-generator/cmd/api"
	"gitlab.r-styleserv.com/r-style/devops/traefik-hosts-generator/cmd/file_writer"
	"gitlab.r-styleserv.com/r-style/devops/traefik-hosts-generator/cmd/generator"
	"log"
	"os"
	"reflect"
	"time"
)

const ApiUrl = "http://localhost:8080/api"
const LocalIP = "127.0.0.1"

func main() {
	apiUrl := flag.String("api", ApiUrl, "specify custom traefik API url, example: 'http://127.0.0.1:8080/api'")
	localIP := flag.String("ip", LocalIP, "specify custom ip to use in hosts file, example: '192.168.33.10'")
	hostsFile := flag.String("file", file_writer.HostsLocation, "specify custom hosts file location, example: '/etc/hosts_custom'")
	watch := flag.Bool("watch", false, "enable API polling mode: true/false")
	platform := flag.String("platform", "", "change line-endings style for hosts file, default: '', available: darwin, windows, linux")
	quiet := flag.Bool("quiet", false, "disable logging")
	period := flag.Int("freq", 5, "poll every N seconds")
	flag.Parse()

	Log := func(fmt string, params ...interface{}) {
		if !*quiet {
			log.Printf(fmt, params)
		}
	}

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
			Log("Unknown platform specified: %s, supported: linux, darwin, windows", *platform)
			os.Exit(1)
		}
	}


	var prevHosts []string
	for {
		hosts, err := api.GetHosts(*apiUrl)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(prevHosts, hosts) {
			fileContent := generator.GenerateStrings(hosts, *localIP, lineEnding)

			err = file_writer.WriteToHosts(fileContent, *hostsFile, lineEnding)
			if err != nil {
				panic(err)
			}

			prevHosts = hosts

			Log("updated hosts file, new hosts: %s%s%s", lineEnding, fileContent, lineEnding)
		} else {
			Log("traefik hosts didn't change, skipping %s", lineEnding)
		}

		if !*watch {
			break
		}

		time.Sleep(time.Second * time.Duration(*period))
	}
}
