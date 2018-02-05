all: install_requirements get_deps build test

install_requirements:
	go get -u github.com/golang/dep/cmd/dep

get_deps:
	dep ensure

build:
	go install

test:
	go test -coverprofile=coverage.out

show_coverage:
	go tool cover -html=coverage.out

.PHONY: install_requirements get_deps build test show_coverage
