############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Build the binary statically.
ENV CGO_ENABLED=0


# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
# Build the binary.
RUN go build -o /app/generator .

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /app/generator /generator
# Run the binary.
ENTRYPOINT ["/generator"]
