FROM golang:alpine AS build
ENV GOPROXY=https://proxy.golang.org

WORKDIR /go/src/portfolio

RUN apk update && apk add --no-cache git

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/portfolio cmd/web/*

RUN ls ui/

EXPOSE 4000

FROM alpine
RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=build /go/bin/portfolio /go/bin/portfolio
COPY --from=build /go/src/portfolio /go/src/portfolio

ENTRYPOINT [ "go/bin/portfolio" ]