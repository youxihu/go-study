GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

.PHONY: init
# init env
init:
	go get entgo.io/ent/entc/gen@v0.14.1
	go get entgo.io/ent/cmd/internal/printer@v0.14.1
	go get entgo.io/ent/cmd/ent@v0.14.1



# generate ent code
ent:
	@go run entgo.io/ent/cmd/ent generate \
				--feature privacy \
				--feature sql/modifier \
				--feature intercept,schema/snapshot \
				--feature entql \
				--feature sql/upsert \
				./internal/data/ent/schema
