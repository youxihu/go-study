.PHONY: build
# build
build:
	rm -rf ./bin
	mkdir -p bin/ && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags "-H windowsgui" -o ./bin/ ./...
	upx  ./bin/*


.PHONY: wire
# wire
wire:
	wire cmd/tabtab/wire.go

.PHONY: init
# init env
init:
	go get entgo.io/ent/entc/gen@v0.14.1
	go get entgo.io/ent/cmd/internal/printer@v0.14.1
	go get entgo.io/ent/cmd/ent@v0.14.1
	sudo apt-get install -y \
        pkg-config \
        libgl1-mesa-dev \
        libx11-dev \
        libxcursor-dev \
        libxrandr-dev \
        libxi-dev \
        libxinerama-dev \
        libxxf86vm-dev \
        libxext-dev

# generate ent code
ent:
	@go run entgo.io/ent/cmd/ent generate \
		--feature privacy \
		--feature sql/modifier \
		--feature intercept,schema/snapshot \
		--feature entql \
		--feature sql/upsert \
		./internal/data/ent/schema
