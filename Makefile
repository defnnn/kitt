SHELL := /bin/bash

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

KITT_IP := 169.254.32.1

clean: # Remove certs
	rm -f etc/acme/acme.json

up: # Bring up networking and kitt
	cp etc/acme/acme.json.whoa.bot etc/acme/acme.json
	chmod 600 etc/acme/acme.json
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		mount bpffs /sys/fs/bpf -t bpf
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "ip link add dummy0 type dummy; ip addr add $(KITT_IP)/32 dev dummy0; ip link set dev dummy0 up"
	docker network create --subnet 172.31.188.0/24 kitt || true
	docker-compose down || trued
	docker-compose up -d

macos:
	for ip in $(KITT_IP); do sudo ifconfig lo0 alias "$$ip" netmask 255.255.255.255; done

down: # Shut down docker-compose and dummy interface
	docker-compose down --remove-orphans || true
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "ip addr del $(KITT_IP)/32 dev dummy0"

	docker rm -f app1 || true
	docker run -d --rm --name app1 --net cnet -l "id=app1" cilium/demo-httpd

demo:
	docker network create --driver cilium --ipam-driver cilium cnet || true
	docker-compose exec -T cilium cilium policy delete --all
	cat demo.yml | docker run --rm -i letfn/python-cli yq . | docker-compose exec -T cilium cilium policy import -
	docker rm -f app1 || true
	docker run -d --rm --name app1 --net cnet -l "id=app1" cilium/demo-httpd
	@figlet good
	-docker run --rm -ti --net cnet -l "id=app2" cilium/demo-client curl http://app1/public
	@figlet no good
	-docker run --rm -ti --net cnet -l "id=app2" cilium/demo-client curl -X POST -d {} http://app1/public
	@figlet dropped
	-docker run --rm -ti --net cnet -l "id=app3" cilium/demo-client curl -m 2 http://app1/public
	@figlet denied
	-docker run --rm -ti --net cnet -l "id=app2" cilium/demo-client curl -si http://app1/private
