FROM golang:latest as builder
LABEL maintainer="grindlemire"

WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install openssl 
RUN apt-get update && \
    apt-get install libssl-dev
# Generate a self signed certificate
RUN openssl req -new -newkey rsa:2048 -days 365 -nodes \
    -subj "/O=grindlemire" \
    -x509 -keyout server.key -out server.crt 

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app


######## Start a new stage from scratch #######
FROM alpine:latest  
WORKDIR /root/
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/app /app/env /app/server.crt /app/server.key ./
EXPOSE 443
EXPOSE 80
CMD ["./app"]
