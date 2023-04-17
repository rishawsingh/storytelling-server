FROM golang:alpine as story-magic-server

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /server
COPY go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -o bin/story-magic cmd/main.go


FROM alpine:latest

WORKDIR /

COPY --from=story-magic-server /server/bin .
COPY --from=story-magic-server /server/database/migrations ./database/migrations

EXPOSE 8092
ENTRYPOINT ["./story-magic"]
