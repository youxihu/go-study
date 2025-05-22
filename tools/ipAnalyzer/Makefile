.PHONY: build
# build
build:
	rm -rf ./bin
	mkdir -p bin/ &&  CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o bin/ipanalyzer cmd/main.go
	upx  ./bin/*


