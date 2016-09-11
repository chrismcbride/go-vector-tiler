PKGS=$(shell go list ./... | grep -v /vendor/)

.PHONY: install_deps
install_deps:
	glide install

.PHONY: fmt
fmt:
	go get golang.org/x/tools/cmd/goimports
	for pkg in $(PKGS) ; do \
		$(GOPATH)/bin/goimports -w $(GOPATH)/src/$$pkg/*.go ; \
	done
	go fmt $(PKGS)

.PHONY: test
test: fmt
	go vet $(PKGS)
	go get github.com/golang/lint/golint
	for pkg in $(PKGS) ; do \
		$(GOPATH)/bin/golint -set_exit_status $$pkg ; \
	done
	go test -cover -race $(PKGS)

.PHONY: build
build: fmt
	go build .
