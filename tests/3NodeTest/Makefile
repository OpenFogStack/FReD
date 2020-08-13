.PHONY: clean start test

clean: ## clean up unused resources
	@docker-compose down --remove-orphan
	@docker network remove fredwork

test: ## start all containers with docker compose
	@docker network create fredwork --internal --gateway 172.26.0.1 --subnet 172.26.0.0/16
	@docker-compose build
	@docker-compose up --force-recreate --abort-on-container-exit --renew-anon-volumes --remove-orphans