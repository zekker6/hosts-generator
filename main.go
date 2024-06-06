package main

import (
	"context"
	"flag"
	logger "log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"hosts-generator/cmd"
	"hosts-generator/cmd/file_writer"
	"hosts-generator/cmd/generator"
	"hosts-generator/cmd/parsers"
	"hosts-generator/cmd/parsers/caddy"
	"hosts-generator/cmd/parsers/kubernetes"
	"hosts-generator/cmd/parsers/traefik_v2"

	"k8s.io/client-go/util/homedir"
)

var (
	localIP   = flag.String("ip", "127.0.0.1", "specify custom ip to use in hosts file, example: '192.168.33.10'")
	hostsFile = flag.String("file", file_writer.HostsLocation, "specify custom hosts file location, example: '/etc/hosts_custom'")
	platform  = flag.String("platform", "", "change line-endings style for hosts file, default: '', available: darwin, windows, linux")
	quiet     = flag.Bool("quiet", false, "disable logging")
	period    = flag.Int("freq", 5, "poll every N seconds")
	watch     = flag.Bool("watch", false, "enable API polling mode: true/false")
	postfix   = flag.String("postfix", "", "use unique postfix if 2 parallel instances are running")

	skipWildcard = flag.Bool("skipWildcard", false, "remove wildcard entries in hosts file. "+
		"Not all DNS servers support wildcard entries, so this option can be used to filter out unsupported entries.")

	kubeConfig = flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "specify full path to kubeconfig")

	traefikUrl = flag.String("traefikUrl", "http://localhost:8080/api", "specify custom traefik API url, example: 'http://127.0.0.1:8080/api'")

	caddyURL = flag.String("caddyURL", "", "specify custom caddy API url, example: 'http://127.0.0.1:2019/config/'")
)

func main() {
	flag.Parse()

	lineEnding := detectLineEndings()

	adapter := file_writer.NewFileHostsAdapter(*hostsFile)

	writer := file_writer.NewWriter(&adapter, lineEnding, *postfix)

	handleExit(writer)

	clients := buildClientsConfig()
	if len(clients) == 0 {
		log("WARN: no clients configured")
	}

	appConfig := &cmd.AppConfig{
		Clients:      clients,
		Writer:       writer,
		LineEnding:   lineEnding,
		TargetIP:     *localIP,
		SyncPeriod:   time.Second * time.Duration(*period),
		EnableWatch:  *watch,
		Logger:       log,
		SkipWildcard: *skipWildcard,
	}
	app := cmd.NewApp(appConfig)

	err := app.Run(context.Background())
	if err != nil {
		log("runtime error: %+v", err)
		err := app.Stop()
		if err != nil {
			log("failed to clear hosts: %+v", err)
		}
		os.Exit(1)
	}
}

func buildClientsConfig() []parsers.Parser {
	type clientConf struct {
		enable bool
		client parsers.Parser
	}

	clientsConf := []clientConf{
		{len(*kubeConfig) != 0, kubernetes.NewKubernetesClient(*kubeConfig)},
		{len(*traefikUrl) != 0, traefik_v2.NewTraefikV2Client(*traefikUrl)},
		{len(*caddyURL) != 0, caddy.NewCaddyV3(*caddyURL)},
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
