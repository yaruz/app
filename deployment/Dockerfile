# README run from root dir (./..)
# docker image build -f deployment/DockerFile -t app:1.0 .
# docker container run --publish 8080:8080 --name app app:1.0

FROM golang:1.18-alpine as builder

WORKDIR /usr/src/app
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -v -o /go/bin/restapi ./cmd/restapi


FROM alpine:latest

#RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY metadata /bin/metadata
COPY --from=builder /go/bin/restapi .

EXPOSE 80

CMD ["restapi"]
