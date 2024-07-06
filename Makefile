default: test

SHELL=/usr/bin/env bash
GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover

.PHONY: test
test: test/cover test/report

.PHONY: test/cover
test/cover:
	$(GOTEST) -v -race -failfast -parallel 4 -cpu 4 -coverprofile main.cover.out ./...

.PHONY: test/report
test/report:
	if [[ "$$HOSTNAME" =~ "codespaces-"* ]]; then \
		mkdir -p /tmp/stringable_test; \
		$(GOCOVER) -html=main.cover.out -o /tmp/stringable_test/coverage.html; \
		sudo python -m http.server -d /tmp/stringable_test -b localhost 80; \
	else \
		$(GOCOVER) -html=main.cover.out; \
	fi