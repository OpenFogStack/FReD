.PHONY: clean fred-debug-nodeB fred failtest alexandra

clean: ## clean up all resources
	@docker-compose -f fredwork.yml -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f 3NodeTester.yml -f trigger.yml down --remove-orphan

fred: clean ## Start all containers except the tester so that other clients can access them
	@docker-compose -f fredwork.yml -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml build
	@docker-compose -f fredwork.yml -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml up --force-recreate --renew-anon-volumes --remove-orphans

fred-debug-nodeB: clean ## same as "fred" but with nodeB in debugger mode
	@docker-compose -f fredwork.yml -f etcd.yml -f nodeA.yml -f nodeB-debug.yml -f nodeC.yml -f trigger.yml build
	@docker-compose -f fredwork.yml -f etcd.yml -f nodeA.yml -f nodeB-debug.yml -f nodeC.yml -f trigger.yml up --force-recreate --renew-anon-volumes --remove-orphans

failtest: ## Start the failtest
	@docker-compose -f failtest.yml build
	@docker-compose -f failtest.yml up --force-recreate --abort-on-container-exit --renew-anon-volumes --exit-code-from tester

alexandra: ## alexandra (starts 3node, alexandra)
	@docker-compose -f fredwork.yml -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml -f alexandra.yml build
	@docker-compose -f fredwork.yml -f etcd.yml -f nodeA.yml -f nodeB.yml -f nodeC.yml -f trigger.yml -f alexandra.yml up --force-recreate --renew-anon-volumes --remove-orphans
