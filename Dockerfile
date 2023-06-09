FROM golang:latest

RUN apt-get update && apt-get install make bash

WORKDIR /app

COPY ./ /app

RUN go mod download -x

COPY --from=itinance/swag /root/swag /usr/local/bin

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go install -mod=mod github.com/swaggo/swag/cmd/swag@latest

ENTRYPOINT CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="make build-dev" --command=./main
