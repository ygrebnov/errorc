ROOT_PATH := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
COVERAGE_PATH := $(ROOT_PATH).coverage/

test:
	@rm -rf $(COVERAGE_PATH)
	@mkdir -p $(COVERAGE_PATH)
	@go test -v -coverpkg=./... ./... -coverprofile $(COVERAGE_PATH)coverage.txt
	@go tool cover -html=$(COVERAGE_PATH)coverage.txt -o $(COVERAGE_PATH)coverage.html

bench:
	@go test -bench=.

fuzz:
	@go test -fuzz=Fuzz -fuzztime=10s ./...

.PHONY: test bench fuzz