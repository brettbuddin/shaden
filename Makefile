default: build

install:
	go install .

clean:
	rm -rf bin/

build:
	go build -o bin/shaden .

test:
	go test -race -cover ./...

vet:
	go vet ./...
 
lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint

ci: test lint vet

.PHONY: build install clean test vet lint ci
