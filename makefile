VersionFile=VERSION
VERSION=`cat $(VersionFile)`

start:
	docker-compose up -d

stop:
	docker-compose stop

build:
	docker-compose build

show-version:
	echo ${VERSION}

inc-version:
	go run cmd/release/inc-version.go

rebuild:
	go build && swag init && make fix-swagger-models

run-with-swagger:
	go build && swag init && make fix-swagger-models && go run main.go

# fix-swagger-models
# swaggo/swag has a bug that will prevent renaming of Models from "model.Account" ino "Account"
# we are going to fix this generation with this command

fix-swagger-models:
	chmod +x fix-swagger-files.sh

build-dev:
	swag init
	chmod +x fix-swagger-files.sh
	go build -v main.go

