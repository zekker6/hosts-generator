OUT_DIR=./out

build: build_linux build_mac build_win

build_linux: prepare_out_dir
	go build -o ${OUT_DIR}/linux.traefik-hosts-generator

build_mac: prepare_out_dir
	GOOS=darwin GOARCH=amd64 go build -o ${OUT_DIR}/darwin.traefik-hosts-generator

build_win: prepare_out_dir
	GOOS=windows GOARCH=amd64 go build -o ${OUT_DIR}/win.traefik-hosts-generator

docker_build: build_linux
	docker build -f docker/Dockerfile ${OUT_DIR} -t zekker6/traefik-hosts-generator

docker_push: docker_build
	docker push zekker6/traefik-hosts-generator

prepare_out_dir:
	mkdir -p ${OUT_DIR}

test:
	go test -race -short `go list ./... | grep -v /vendor/`