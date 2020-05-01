build:
	cd src && go build
	
distribute:
	./environment/build-scripts/distribute-executables.sh

fix-permissions:
	chmod u+x ./environment/build-scripts/install-dependencies.sh
	chmod u+x ./environment/build-scripts/distribute-executables.sh

help:
	cd src && go run main.go -h

install:
	./environment/build-scripts/install-dependencies.sh

build-local:
	docker build -t gofire .

run-local-docker:
	docker run \
		-i -t \
		--rm \
		-p 8000:8000 \
		-e GOFIRE_MASTER=true \
		-e GOFIRE_MASTER_HOST= \
		gofire:latest
	

run-master:
	cd src && GOFIRE_MASTER=true go run main.go

run-slave:
	cd src && go run main.go

run-slave2:
	cd src && GOFIRE_PORT=8002 GO111MODULE=on go run main.go

run-js:
	cd frontend && npm run start

build-js:
	cd frontend && npm run build
  
