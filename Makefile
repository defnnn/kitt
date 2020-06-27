SHELL := /bin/bash

.PHONY: test

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

test: # Run tests
	cd test && $(MAKE)
	cd test && git diff

drone-test: # Run tests with drone specific setup
	mkdir /tmp/src
	rsync -ia . /tmp/src/.
	cd /tmp/src/test && $(MAKE)
	cd /tmp/src/test && git diff

setup once:
	exec/kitt-setup

dc0:
	exec/kitt-consul-run

dc0-gateway:
	. .env; export CONSUL_HTTP_TOKEN; env CONSUL_HTTP_ADDR=169.254.32.10:8500 CONSUL_GRPC_ADDR=169.254.32.10:8502 consul connect envoy -mesh-gateway -register -address 169.254.32.10:4433

dc0-proxy:
	. .env; export CONSUL_HTTP_TOKEN; env CONSUL_HTTP_ADDR=169.254.32.10:8500 CONSUL_GRPC_ADDR=169.254.32.10:8502 consul connect envoy -sidecar-for dc0 -admin-bind 127.0.0.1:19002

dc0-test:
	@curl http://localhost:9091
