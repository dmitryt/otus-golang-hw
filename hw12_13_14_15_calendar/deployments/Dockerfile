FROM golang:1.14-alpine AS builder

RUN apk update

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

ARG MAIN_FILE_PATH

# Build the application
RUN go build -o main $MAIN_FILE_PATH
RUN chmod +x main

# Build a small image
FROM scratch

COPY --from=builder /build/main /
EXPOSE 8081

# Command to run
ENTRYPOINT ["/main"]