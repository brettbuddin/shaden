test:
	go test -race -cover ./...

vet:
	go vet ./...
 
lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint

ci: test lint vet

.PHONY: test vet lint ci
