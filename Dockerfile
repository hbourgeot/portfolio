FROM golang

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd/ ./
COPY ./internal/models/forms.go ./

RUN go build cmd/web

EXPOSE 4000



CMD ["/server"]