# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# Stage 2: Create the runtime image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

ENV APPLICATION_OPENAI_KEY=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
ENV APPLICATION_WHISPER_URL=https://api.openai.com/v1/audio/transcriptions
ENV APPLICATION_GPT_URL=https://api.openai.com/v1/chat/completions
ENV PORT=8089

EXPOSE 8089

# Command to run the executable
CMD ["./main"]