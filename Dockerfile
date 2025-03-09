FROM golang:1.22-alpine3.20
LABEL authors="itzikreuven"

COPY ./ /app
WORKDIR /app/backend
RUN apk update
RUN apk add curl
RUN go mod download && go mod verify
RUN go build .

ENTRYPOINT ["./backend"]