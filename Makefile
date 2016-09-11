PKGS=$(shell go list ./... | grep -v /vendor/)

.PHONY: build
build:
	go vet $(PKGS)
	go get github.com/golang/lint/golint
	for pkg in $(PKGS) ; do \
		$(GOPATH)/bin/golint $$pkg ; \
	done
	go test -v -cover -race $(PKGS)
	go build .

.PHONY: fmt
fmt:
	go get golang.org/x/tools/cmd/goimports
	for pkg in $(PKGS) ; do \
		$(GOPATH)/bin/goimports -w $(GOPATH)/src/$$pkg/*.go ; \
	done
	go fmt $(PKGS)
