# Hosts generator

A small tool which is able to generate hosts file content for services discovered from different sources:
* [Traefik](https://traefik.io)
* Kubernetes
* [Caddy](https://caddyserver.com/)

Available as docker image:
- v0.6.0+ published [here](https://github.com/zekker6/hosts-generator/pkgs/container/hosts-generator)
- previous versions published [here](https://github.com/users/zekker6/packages/container/package/traefik-hosts-generator)

# Usage

Initially it was developed to enhance local development toolchain which helps to work on several projects at the same time. Workflow suggested having multiple docker-compose apps running and having Traefik as reverse-proxy with dynamic discovery over docker socket.

Example config:
```yaml

version: "3"

services:
  traefik:
    image: traefik:v3.0.1
    restart: unless-stopped
    labels:
      traefik.port: 8080
    volumes:
      - "./traefik.toml:/etc/traefik/traefik.toml"
      - "/var/run/docker.sock:/var/run/docker.sock"
    ports:
      - "80:80"
    networks:
      - tk_web

  tk-hosts:
    image: ghcr.io/zekker6/hosts-generator:v1.0.1
    restart: unless-stopped
    volumes:
      - /etc/hosts:/hosts
    command: "-platform=linux -traefikUrl=http://traefik:8080/api -file=/hosts -watch=true -freq=10"
    networks:
      - tk_web
    depends_on:
      - traefik

networks:
  tk_web:
    external: true
```

This config will use external network `tk_web` to establish communication with application container.
It is needed to add this network to all containers which will be connected via Traefik.

You can create this network with the following command:

```sh
docker network create tk_web
```

Fully working example can be found at [examples folder](example/).

## Command line flags

CLI flags will allow to override default behaviour such as line endings for different host operating systems(useful when using docker image), changing generated block postfix(to allow using several concurrent instances of traefik generator).

```
  -caddyURL string
        specify custom caddy API url, example: 'http://127.0.0.1:2019/config/'
  -file string
        specify custom hosts file location, example: '/etc/hosts_custom' (default "/etc/hosts")
  -freq int
        poll every N seconds (default 5)
  -ip string
        specify custom ip to use in hosts file, example: '192.168.33.10' (default "127.0.0.1")
  -kubeconfig string
        specify full path to kubeconfig (default "/home/zekker/.kube/config")
  -platform string
        change line-endings style for hosts file, default: '', available: darwin, windows, linux
  -postfix string
        use unique postfix if 2 parallel instances are running
  -quiet
        disable logging
  -skipWildcard
        remove wildcard entries in hosts file. Not all DNS servers support wildcard entries, so this option can be used to filter out unsupported entries.
  -traefikUrl string
        specify custom traefik API url, example: 'http://127.0.0.1:8080/api' (default "http://localhost:8080/api")
  -watch
        enable API polling mode: true/false
```
