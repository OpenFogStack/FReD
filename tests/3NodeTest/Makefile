.PHONY: clean start test

clean: ## clean up unused resources
	@docker-compose -f etcd.yml -f stores.yml -f nodes.yml -f tester.yml down --remove-orphan
	@docker network remove fredwork

test: ## start all containers with docker compose and run the test dockerized
	@docker network create fredwork --internal --gateway 172.26.0.1 --subnet 172.26.0.0/16
	@docker-compose -f etcd.yml -f stores.yml -f nodes.yml -f tester.yml build
	@docker-compose -f etcd.yml -f stores.yml -f nodes.yml -f tester.yml up --force-recreate --abort-on-container-exit --renew-anon-volumes --remove-orphans