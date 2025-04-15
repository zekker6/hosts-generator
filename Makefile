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

tests_integration_prepare:
	docker network create tk_web_tests || true
	docker-compose -f ./tests/integration/docker-compose.yml -p tk-hosts-tests up -d

tests_integration_stop:
	docker-compose -f ./tests/integration/docker-compose.yml -p tk-hosts-tests down
	docker network rm tk_web_tests

test_integration: tests_integration_prepare
	go test -tags integration ./tests/integration/

lint: lint_vet lint_gofmt lint_golangcilint

lint_vet:
	go vet ./...

lint_gofmt:
	gofmt -l -w -s ./

lint_golangcilint:
	golangci-lint run --modules-download-mode readonly --timeout=10m
