COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')
PWD=$(shell pwd)
all: build

build: ztserver
web:
	cd dashboard && npm install && npm run build:prod
	mkdir -p server/dashboard
	cp -rf dashboard/dist server/dashboard
ztserver:web
	cd $(PWD)/server && go build -ldflags "-X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\"" -o ./bin/ztserver ./cmd/ztserver
run:ztserver
	cd $(PWD)/server && ./bin/ztserver -config etc/config.yaml
clean:
	cd $(PWD)
	@rm -rf bin && rm -rf dashboard/dist && rm -rf server/bin/* && rm -rf server/dashboard
