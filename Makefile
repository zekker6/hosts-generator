build:
	go build -o out/traefik-hosts-generator

docker_build: build
	docker build -f docker/Dockerfile out/ -t zekker6/traefik-hosts-generator

docker_push: docker_build
	docker push zekker6/traefik-hosts-generator