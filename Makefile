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
	consul agent -config-file="$(PWD)/etc/consul_config/dc0.hcl" -data-dir="$(PWD)/etc/consul_dc0" -join-wan=169.254.32.1

dc0-gateway:
	set -a; . .env; set +a; exec/kitt-dc0 connect envoy -mesh-gateway -register -address 169.254.32.0:4444

dc0-proxy:
	set -a; . .env; set +a; exec/kitt-dc0 connect envoy -sidecar-for dc0 -admin-bind 127.0.0.1:19001

dc0-test:
	@curl http://localhost:9091
