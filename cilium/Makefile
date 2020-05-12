SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

up: # Bring up networking and kitt
	docker run --rm -i --privileged --network=host --pid=host alpine nsenter -t 1 -m -u -n -i -- \
		mount bpffs /sys/fs/bpf -t bpf
	docker-compose down || true
	docker-compose up -d

network:
	docker network create --subnet 172.31.189.0/24 --driver cilium --ipam-driver cilium cnet

down: # Shut down docker-compose and dummy interface
	docker-compose down --remove-orphans || true

demo:
	docker-compose exec -T cilium cilium policy delete --all
	cat demo.yml | docker run --rm -i letfn/python-cli yq . | docker-compose exec -T cilium cilium policy import -
	docker rm -f app1 || true
	docker run -d --rm --name app1 --net kitt -l "id=app1" cilium/demo-httpd
	@figlet good
	-docker run --rm -ti --net kitt -l "id=app2" cilium/demo-client curl http://app1/public
	@figlet no good
	-docker run --rm -ti --net kitt -l "id=app2" cilium/demo-client curl -X POST -d {} http://app1/public
	@figlet dropped
	-docker run --rm -ti --net kitt -l "id=app3" cilium/demo-client curl -m 2 http://app1/public
	@figlet denied
	-docker run --rm -ti --net kitt -l "id=app2" cilium/demo-client curl -si http://app1/private
