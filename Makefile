build:
	go build -o kgrep main.go

install: build
	mv kgrep $(GOPATH)/bin

clean:
	rm -fr kgrep

test:
	go test ./...
