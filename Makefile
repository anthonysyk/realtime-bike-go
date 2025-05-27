clean:
	@echo "> cleaning..."
	docker-compose down

init:
	@echo "> initialization..."
	docker-compose up -d

	@echo "> waiting for mongodb replication to be ready (this can take some seconds...)"
	-for i in {1..5}; do docker-compose exec mongo-primary mongo \
		'mongodb://root:secret@127.0.0.1:27011,127.0.0.1:27012,127.0.0.1:27013/watcher?connect=replicaset&replicaSet=replicaset&authSource=admin' \
		--eval 'rs.status();' > /dev/null && break || sleep 1; done


run-collector:
	@echo "> launching collector..."
	go run cmd/collector/main.go

run-watcher:
	@echo "> launching watcher..."
	go run cmd/watcher/main.go

.PHONY: clean init