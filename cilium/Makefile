SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

setup:
	$(MAKE) once
	docker-compose down --remove-orphans
	docker-compose up -d --remove-orphans

once:
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		mount bpffs /sys/fs/bpf -t bpf

status:
	docker-compose exec cilium cilium status --all-controllers --all-health --all-redirects
	docker-compose exec cilium hubble status

observe:
	@docker-compose exec cilium hubble observe -f -j
