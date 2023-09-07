## Run integration tests
integration:
	docker-compose -f test/integration/docker-compose.yml up --abort-on-container-exit --exit-code-from runner

## Cleanup
clean:
	docker-compose -f test/integration/docker-compose.yml rm
