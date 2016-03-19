#! /usr/bin/make
#
# Makefile for encoding
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "lint" runs the linter and checks the code format using goimports
# - "test" runs the tests
#
# Meta targets:
# - "all" is the default target, it runs all the targets in the order above.
#http://goa.design/reference/goa/design/apidsl.html#func-consumes-a-name-apidsl-consumes-a:aab4f9d6f98ed71f45bd470427dde2a7
DIRS=$(shell go list -f {{.Dir}} ./...)
DEPEND=\
	github.com/asaskevich/govalidator \
	github.com/goadesign/goa.design/tools/mdc \
	github.com/goadesign/goa.design/tools/godoc2md \
	github.com/golang/lint/golint \
	github.com/on99/gocyclo \
	github.com/onsi/ginkgo \
	github.com/onsi/ginkgo/ginkgo \
	github.com/onsi/gomega \
	github.com/spf13/hugo \
	golang.org/x/tools/cmd/cover \
	gopkg.in/yaml.v2

all: depend lint cyclo test

docs:
	@git clone https://github.com/goadesign/goa.design
	@rm -rf goa.design/content/reference goa.design/public
	@mdc github.com/goadesign/encoding goa.design/content/reference --exclude goa.design
	@cd goa.design && hugo
	@rm -rf public
	@mv goa.design/public public
	@rm -rf goa.design

depend:
	@go get $(DEPEND)

lint:
	@for d in $(DIRS) ; do \
		if [ "`goimports -l $$d/*.go | tee /dev/stderr`" ]; then \
			echo "^ - Repo contains improperly formatted go files" && echo && exit 1; \
		fi \
	done
	@if [ "`golint ./... |  tee /dev/stderr`" ]; then \
		echo "^ - Lint errors!" && echo && exit 1; \
	fi

cyclo:
	@if [ "`gocyclo -over 20 . | grep -v examples/cellar | tee /dev/stderr`" ]; then \
		echo "^ - Cyclomatic complexity exceeds 20, refactor the code!" && echo && exit 1; \
	fi

test:
	@ginkgo -r --randomizeAllSpecs --failOnPending --randomizeSuites --race -skipPackage vendor
