# Specify the base image for the go app.
FROM golang:alpine
# Specify that we now need to execute any commands in this directory.
WORKDIR /go/src/github.com/henrry.online

RUN apk add --no-cache go

# Copy everything from this project into the filesystem of the container.
COPY . .
# Obtain the package needed to run redis commands. Alternatively use GO Modules.
RUN go mod download

# Compile the binary exe for our app.
RUN go build -o main .
# Start the application.
CMD ["./main"]