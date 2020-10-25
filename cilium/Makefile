SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

setup:
	$(MAKE) once up

dummy:
	docker run --rm -i --privileged --network=host --pid=host alpine \
		nsenter -t 1 -m -u -n -i -- bash -c "ip link add dummy0 type dummy; ip addr add 169.254.32.1/32 dev dummy0; ip link set dev dummy0 up"

logs:
	docker-compose logs -f

recreate:
	$(MAKE) down up

down:
	docker-compose down --remove-orphans

up:
	docker-compose up -d --remove-orphans

once:
	docker run --rm -i --privileged --network=host --pid=host cilium/cilium \
		iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
	docker run --rm -i --privileged --network=host --pid=host alpine \
		nsenter -t 1 -m -u -n -i -- mount bpffs /sys/fs/bpf -t bpf

policy:
	docker-compose exec -T cilium cilium policy delete --all
	cat policy.yml | docker run --rm -i letfn/python-cli yq . | docker-compose exec -T cilium cilium policy import -

status:
	@docker-compose exec cilium cilium status --all-controllers --all-health --all-redirects; echo
	@docker-compose exec cilium hubble status; echo
	@docker-compose exec cilium cilium node list; echo

observe:
	@docker-compose exec cilium hubble observe -f -j --verdict FORWARDED \
		|	while read -r a; do echo "$$a" | jq . | perl -pe 's{$$}{\r}'; echo; echo; echo; done

dropped:
	@docker-compose exec cilium hubble observe -f -j --verdict DROPPED --protocol tcp \
		|	while read -r a; do echo "$$a" | jq . | perl -pe 's{$$}{\r}'; echo; echo; echo; done

dropped-udp:
	@docker-compose exec cilium hubble observe -f -j --verdict DROPPED --protocol udp\
		|	while read -r a; do echo "$$a" | jq . | perl -pe 's{$$}{\r}'; echo; echo; echo; done
