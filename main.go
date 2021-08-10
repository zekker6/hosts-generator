package main

import (
	"flag"
	"hosts-generator/cmd/file_writer"
	"hosts-generator/cmd/generator"
	"hosts-generator/cmd/parsers"
	"hosts-generator/cmd/parsers/kubernetes"
	"hosts-generator/cmd/parsers/traefik"
	"k8s.io/client-go/util/homedir"
	logger "log"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"syscall"
	"time"
)

var (
	localIP   = flag.String("ip", "127.0.0.1", "specify custom ip to use in hosts file, example: '192.168.33.10'")
	hostsFile = flag.String("file", file_writer.HostsLocation, "specify custom hosts file location, example: '/etc/hosts_custom'")
	platform  = flag.String("platform", "", "change line-endings style for hosts file, default: '', available: darwin, windows, linux")
	quiet     = flag.Bool("quiet", false, "disable logging")
	period    = flag.Int("freq", 5, "poll every N seconds")
	watch     = flag.Bool("watch", false, "enable API polling mode: true/false")
	postfix   = flag.String("postfix", "", "use unique postfix if 2 parallel instances are running")

	kubeConfig = flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "specify full path to kubeconfig")
	kubeEnable = flag.Bool("kube", false, "enable kube client")

	traefikProvider = flag.String("traefikProvider", "docker", "traefik traefikProvider to use")
	traefikUrl      = flag.String("traefikUrl", "http://localhost:8080/api", "specify custom traefik API url, example: 'http://127.0.0.1:8080/api'")
	traefikEnable   = flag.Bool("traefik", false, "enable traefik client")
)

func main() {
	flag.Parse()

	lineEnding := detectLineEndings()

	adapter := file_writer.NewFileHostsAdapter(*hostsFile)

	writer := file_writer.NewWriter(&adapter, lineEnding, *postfix)

	defer func() {
		log("clearing before exit")
		err := writer.Clear()
		if err != nil {
			log("failed to clear hosts: %+v", err)
		}
	}()

	handleExit(writer)

	var prevHosts []string

	clients := buildClientsConfig()

	for {
		hosts := getHosts(clients)

		if !reflect.DeepEqual(prevHosts, hosts) {
			err := writeHosts(hosts, lineEnding, writer)
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

func getHosts(clients []parsers.Parser) []string {
	hosts := make([]string, 0)

	for _, c := range clients {

		clientHosts, err := c.Get()
		if err != nil {
			panic(err)
		}

		hosts = append(hosts, clientHosts...)
	}
	return hosts
}

func buildClientsConfig() []parsers.Parser {
	type clientConf struct {
		enable bool
		client parsers.Parser
	}

	clientsConf := []clientConf{
		{*kubeEnable, kubernetes.NewKubernetesClient(*kubeConfig)},
		{*traefikEnable, traefik.NewTraefikV1Client(*traefikUrl, *traefikProvider)},
	}

	clients := make([]parsers.Parser, 0)

	for _, cc := range clientsConf {
		if cc.enable {
			clients = append(clients, cc.client)
		}
	}

	logger.Println("loaded clients", len(clients))
	return clients
}

func handleExit(writer *file_writer.Writer) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log("stop signal received")
		err := writer.Clear()
		if err != nil {
			log("failed to clear hosts: %+v", err)
		}
		os.Exit(0)
	}()

}

func writeHosts(hosts []string, lineEnding string, writer *file_writer.Writer) error {
	fileContent := generator.GenerateStrings(hosts, *localIP, lineEnding)

	return writer.WriteToHosts(fileContent)
}

func log(fmt string, params ...interface{}) {
	if *quiet {
		return
	}

	if len(params) == 0 {
		logger.Printf(fmt)
	} else {
		logger.Printf(fmt, params...)
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
