# Use an official Golang runtime as a parent image
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace.
COPY ./cmd .

# Build your program for multiple architectures
RUN CGO_ENABLED=0 GOOS=linux go build -v -o website_monitor

# Use a Docker multi-stage build to create a lean production image.
# (Note: 'scratch' gives you a very minimal environment)
FROM scratch

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/app/website_monitor /website_monitor

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/website_monitor"]

# Document that the service listens on port 8080.
# EXPOSE 8080
