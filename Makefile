.PHONY: default
default: build

.PHONY: build
build:
	go build .

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench . ./... -benchmem -run=^#

.PHONY: lint
lint:
	golangci-lint run

.PHONY: coverage-html
coverage-html: coverage
	go tool cover -html=coverage.out

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out -coverpkg=github.com/leonhfr/honeybadger/... ./...

.PHONY: doc
doc:
	godoc -http=:6060

.PHONY: release
release:
	goreleaser release --snapshot --rm-dist
