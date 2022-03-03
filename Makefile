mongo:
	docker run --name mongo5 --network  simplebank_default -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=root -e MONGO_INITDB_DATABASE=mongodb_simple -d mongo:5.0.6-focal

test:
	go test -v -cover ./...

.PHONY: mongo test