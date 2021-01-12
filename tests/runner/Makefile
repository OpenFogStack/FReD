.PHONY: test clean

clean: ## clean up all resources
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f tester.yml -f trigger.yml down --remove-orphan
	@docker network remove fredwork

test: ## start all containers with docker compose and run the test dockerized
	@docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16 || true
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f tester.yml -f trigger.yml build
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f tester.yml -f trigger.yml up --force-recreate --abort-on-container-exit --renew-anon-volumes --remove-orphans

debug-nodeB: ## same as "test" but with nodeB in debugger mode
	@docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16 || true
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB-debug.yml -f nodeC.yml -f tester.yml -f trigger.yml build
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB-debug.yml -f nodeC.yml -f tester.yml -f trigger.yml up --force-recreate --abort-on-container-exit --renew-anon-volumes --remove-orphans

fred: ## Start all containers except the tester so that other clients can access them
	@docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16 || true
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml build
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml up --force-recreate --renew-anon-volumes --remove-orphans

failtest: ## Start the failtest
	@docker-compose -f failtest.yml build
	@docker-compose -f failtest.yml up --force-recreate --abort-on-container-exit --renew-anon-volumes