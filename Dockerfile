FROM golang:alpine AS build
ENV GOPROXY=https://proxy.golang.org

WORKDIR /go/src/portfolio

COPY . .
RUN GOOS=linux go build -o /go/bin/portfolio cmd/web/*

FROM alpine
COPY --from=build /go/bin/portfolio /go/bin/portfolio
ENTRYPOINT [ "go/bin/portfolio" ]