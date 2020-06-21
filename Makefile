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

KITT_IP := 169.254.32.1

once:
	if [[ "$(shell docker network inspect kitt | jq -r 'map(select(.Name == "kitt")) | length')" == 0 ]]; then $(MAKE) network; fi
	$(MAKE) os-$(shell uname -s)-check

network:
	docker network create --subnet 172.31.188.0/24 kitt

os-Linux-check:
	make os-$(shell uname -s)-up || true

os-Linux-up:
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		bash -c "ip link add dummy0 type dummy; ip addr add $(KITT_IP)/32 dev dummy0; ip link set dev dummy0 up"

os-Linux-down:
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
    bash -c "ip addr del $(KITT_IP)/32 dev dummy0"

os-Darwin-check:
	if ! "$(shell ifconfig lo0 | grep "inet $KITT_IP netmask")"; then make os-$(shell uname -s)-up; fi

os-Darwin-up:
	for ip in $(KITT_IP); do sudo ifconfig lo0 alias "$$ip" netmask 255.255.255.255; done

os-Darwin-down:
	for ip in $(KITT_IP); do sudo ifconfig lo0 -alias "$$ip" netmask 255.255.255.255; done
