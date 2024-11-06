# Use a multi-stage build to keep the final image clean and small
FROM golang:alpine3.19 AS binarybuilder

# Install necessary tools and libraries
RUN apk --no-cache add build-base git linux-pam-dev

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules and sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the binary with CGO disabled for a statically linked executable
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Set the working directory in the container
WORKDIR /app

# Copy the compiled binary and configuration files from the builder stage
COPY --from=binarybuilder /app/main .
COPY --from=binarybuilder /app/config ./config
# Ensure the .gitignore file is copied to the final image if it's required at runtime
COPY --from=binarybuilder /app/.gitignore .


EXPOSE 50051

# Command to run the executable
CMD ["./main"]
