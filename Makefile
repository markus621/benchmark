RD=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

run:
	@go run cmd/benchmark/main.go --config=$(RD)/configs/config.yaml.ctmpl