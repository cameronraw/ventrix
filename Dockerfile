# Start from the official golang image
FROM golang:1.19

# Add Maintainer Info
LABEL maintainer="Cameron Raw <cameron.raw89@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

ENV MONTECRISTO_PORT=8080
ENV MONTECRISTO_ENV=development
ENV MONTECRISTO_KEY=secret
ENV MONTECRISTO_SQL_DSN=
ENV MONTECRISTO_REDIS_URL=127.0.0.1:6379
ENV MONTECRISTO_USE_IN_MEMORY=true

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

