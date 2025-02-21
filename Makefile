build:
	go build -o kgrep main.go

install: build
	mv kgrep $(GOPATH)/bin

clean:
	rm -fr kgrep

test:
	go test ./...

.PHONY: release-snapshot
release-snapshot:
	goreleaser release --snapshot --clean

.PHONY: release-check
release-check:
	goreleaser check
