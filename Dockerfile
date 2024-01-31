# Start from the latest golang base image
FROM golang:alpine3.19 AS binary

# Install git
RUN apk add --no-cache git

# Disable Go's module proxy and checksum database
ENV GOPROXY=direct
ENV GONOSUMDB=*

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod tidy

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Copy the migrations directory
COPY pkg/migrations /app/migrations


WORKDIR /app/cmd/0xbase/

# Debug: Build the Go app with verbose output
RUN go build -v -o ../../main-out .

# Start a new stage from scratch
FROM gcr.io/distroless/base-debian12 AS build-release-stage

# Copy the binary from the previous stage
COPY --from=binary /app/main-out /app/

# Copy the environment file
COPY --from=binary /app/app.env /

# Copy the migrations directory to the final stage
COPY --from=binary /app/migrations /app/migrations

EXPOSE 5050

# Command to run the executable
CMD ["/app/main-out"]