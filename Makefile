SHELL := /bin/bash
KITT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: docs

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

all: # Run everything except build
	$(MAKE) fmt
	$(MAKE) lint
	$(MAKE) docs

fmt: # Format drone fmt
	@echo
	drone exec --pipeline $@

lint: # Run drone lint
	@echo
	drone exec --pipeline $@

docs: # Build docs
	@echo
	drone exec --pipeline $@

build: # Build container
	@echo
	drone exec --pipeline $@

edit:
	docker-compose -f docker-compose.docs.yml up --quiet-pull

requirements:
	@echo
	drone exec --pipeline $@

pkg:
	@echo
	drone exec --pipeline $@

KITT_IP := 169.254.32.1

main: target/main.js # Compile release
	true

clean:
	rm -rf target .shadow-cljs node_modules main
	rm -f etc/acme/acme.json

install:
	npm install

watch:
	./node_modules/.bin/shadow-cljs watch app

repl:
	./node_modules/.bin/shadow-cljs cljs-repl app

release: ./node_modules/.bin/shadow-cljs
	./node_modules/.bin/shadow-cljs release app

node_modules/.bin/shadow-cljs:
	$(MAKE) install

target/main.js:
	$(MAKE) release

up: # Bring up networking and kitt
	cp etc/acme/acme.json.whoa.bot etc/acme/acme.json
	chmod 600 etc/acme/acme.json
	docker-compose down || true
	docker-compose up -d

down: # Shut down docker-compose and dummy interface
	docker-compose down --remove-orphans || true

once:
	$(MAKE) network || true
	$(MAKE) linux|| true
	$(MAKE) macos || true

network:
	docker network create --subnet 172.31.188.0/24 kitt

linux:
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "ip link add dummy0 type dummy; ip addr add $(KITT_IP)/32 dev dummy0; ip link set dev dummy0 up"
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "mkdir -p /etc/systemd/resolved.conf.d"
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "cat $(KITT_DIR)/etc/dns/systemd-resolved/kitt_domain.conf | sed 's/KITT_DOMAIN/$(KITT_DOMAIN)/' > /etc/systemd/resolved.conf.d/kitt_domain.conf"
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "systemctl restart systemd-resolved.service"

linux-down:
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "ip addr del $(KITT_IP)/32 dev dummy0 || true"
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "rm  /etc/systemd/resolved.conf.d/kitt_domain.conf || true"
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "systemctl restart systemd-resolved.service"

macos:
	for ip in $(KITT_IP); do sudo ifconfig lo0 alias "$$ip" netmask 255.255.255.255; done

macos-down:
	for ip in $(KITT_IP); do sudo ifconfig lo0 -alias "$$ip" netmask 255.255.255.255; done
