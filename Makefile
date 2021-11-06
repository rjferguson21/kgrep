build:
	go build -o kpretty main.go

install: build
	mv kpretty $(GOPATH)/bin

clean:
	rm -fr kpretty