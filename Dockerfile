# STAGE 1: Build
FROM golang:1.16.3-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/pricescraper-srv main.go

# STAGE 2: Deployment
FROM alpine:3.13.5

COPY --from=build /go/bin/pricescraper-srv /pricescraper-srv

CMD ["/pricescraper-srv"]
