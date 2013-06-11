
all: butler/version.go build

butler/version.go: VERSION
	bash script/print_version_go.sh > butler/version.go

build:
	GOARCH=386   GOOS=darwin go get -v './butler'
	GOARCH=amd64 GOOS=darwin go get -v './butler'
	GOARCH=386   GOOS=linux  go get -v './butler'
	GOARCH=amd64 GOOS=linux  go get -v './butler'

.PHONY: all build
