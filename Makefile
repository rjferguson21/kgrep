build:
	go build -o kgrep main.go

install: build
	mv kpretty $(GOPATH)/bin

clean:
	rm -fr kpretty