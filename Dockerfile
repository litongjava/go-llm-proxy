# Use the offical Golang image to build the app: https://hub.docker.com/_/golang
FROM litongjava/hertz:0.9.0-golang1.22.1 as builder

# Copy code to the image
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app.out main.go

# Start a new image for production without build dependencies
FROM litongjava/alpine:base

# Copy the app binary from the builder to the production image
COPY --from=builder /app.out /app

# Run the app when the vm starts
CMD ["/app"]