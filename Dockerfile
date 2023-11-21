FROM golang:1.21.4

WORKDIR /usr/src/fumies

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY api ./api
COPY cmd ./cmd
COPY middleware ./middleware
COPY migrations ./migrations
COPY .env ./.env

EXPOSE 8080
EXPOSE 5432

CMD [ "go", "run", "cmd/main.go" ]