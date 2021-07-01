.PHONY: all clean 3nodetest 3nodetest-debug-nodeB 3n-all 3n-standard 3n-replica 3n-trigger 3n-immutable 3n-expiry 3n-selfreplica 3n-auth 3n-concurrency fred failtest alexandra

all: clean 3n-all clean

clean: ## clean up all resources
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f 3NodeTester.yml -f trigger.yml down --remove-orphan
	@docker network remove fredwork

3n-all: export TEST_RANGE = -
3n-standard: export TEST_RANGE = 1-1
3n-replica: export TEST_RANGE = 2-2
3n-trigger: export TEST_RANGE = 3-3
3n-immutable: export TEST_RANGE = 4-4
3n-expiry: export TEST_RANGE = 5-5
3n-selfreplica: export TEST_RANGE = 6-6
3n-auth: export TEST_RANGE = 7-7
3n-concurrency: export TEST_RANGE = 8-8
3n-all 3n-standard 3n-replica 3n-trigger 3n-immutable 3n-expiry 3n-selfreplica 3n-auth 3n-concurrency: 3nodetest

3nodetest: clean ## start all containers with docker compose and run the test dockerized
	@docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16 || true
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f 3NodeTester.yml -f trigger.yml build
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f 3NodeTester.yml -f trigger.yml up --force-recreate --abort-on-container-exit --renew-anon-volumes --remove-orphans

3nodetest-debug-nodeB: export TEST_RANGE = -
3nodetest-debug-nodeB: ## same as "test" but with nodeB in debugger mode
	@docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16 || true
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB-debug.yml -f nodeC.yml -f 3NodeTester.yml -f trigger.yml build
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB-debug.yml -f nodeC.yml -f 3NodeTester.yml -f trigger.yml up --force-recreate --abort-on-container-exit --renew-anon-volumes --remove-orphans

fred: clean ## Start all containers except the tester so that other clients can access them
	@docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16 || true
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml build
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml up --force-recreate --renew-anon-volumes --remove-orphans

failtest: ## Start the failtest
	@docker-compose -f failtest.yml build
	@docker-compose -f failtest.yml up --force-recreate --abort-on-container-exit --renew-anon-volumes

alexandra: clean ## start fred with alexandra
	@docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16 || true
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml -f alexandra.yml build
	@docker-compose -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml -f alexandra.yml up --abort-on-container-exit --force-recreate --renew-anon-volumes --remove-orphans
