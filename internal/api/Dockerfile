FROM golang:alpine AS build

ENV GOPROXY=https://proxy.golang.org

WORKDIR /go/src/hangman
RUN apk update && apk add --no-cache git

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/hangman cmd/web/*

EXPOSE 4000

FROM alpine
RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=build /go/bin/hangman /go/bin/hangman
COPY --from=build /go/src/hangman /go/src/hangman

ENTRYPOINT [ "go/bin/hangman" ]