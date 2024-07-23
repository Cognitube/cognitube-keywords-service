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
ENV APPLICATION_AZURE_KEY=f64dbdf234074181a053e7d8b227067e
ENV APPLICATION_AZURE_URL="https://eastus.api.cognitive.microsoft.com/speechtotext/v3.1/transcriptions"
ENV APPLICATION_CALLBACK_URL=https://moccasin-known-doe.ngrok-free.app/api/v1/callback
ENV APPLICATION_DEBUG=false
ENV APPLICATION_GPT_URL=https://api.openai.com/v1/chat/completions
ENV APPLICATION_KAFKA_HOST=localhost
ENV APPLICATION_KAFKA_PORT=9092
ENV APPLICATION_KAFKA_TOPIC=test
ENV APPLICATION_OPENAI_KEY=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
ENV APPLICATION_WHISPER_URL=https://api.openai.com/v1/audio/transcriptions
ENV AZURE_BLOB_CONNECTION_STRING="DefaultEndpointsProtocol=https;AccountName=cognitube;AccountKey=a1XDmr4IlO9I/tcsuh1akTaGFgmp+nQEoQdA8SlFpmmn7Zi0HKeDMk3ntxWjGI/HMFpQjzBys2ZX+AStqYVfsg==;EndpointSuffix=core.windows.net"
ENV KAFKA_EVENTHUB_NAMESPACE=cognitube-kafka
ENV KAFKA_EVENTHUB_NAME=cognitube
ENV KAFKA_EVENTHUB_CONNECTION_STRING="Endpoint=sb://cognitube-kafka.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=l9+PMVbv8R4LuCtQlPo5x8PIE8jZqn8O4+AEhEMQoqA="

EXPOSE 8089

# Command to run the executable
CMD ["./main"]