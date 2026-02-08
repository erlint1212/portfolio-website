# STAGE 1: The Builder (The Forge)
FROM golang:1.25.5 AS builder

WORKDIR /app

# 1. Install Templ (Tools first)
RUN go install github.com/a-h/templ/cmd/templ@v0.3.977 

# 2. Cache Dependencies (Layers matter!)
COPY go.mod go.sum ./
RUN go mod download

# 3. Copy Source & Generate
COPY . .
RUN ["templ", "generate"]

# 4. Test & Build
RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio ./cmd/web/main.go

# STAGE 2: The Runner (The Vessel)
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy only the binary from the Builder
COPY --from=builder /app/portfolio /portfolio

EXPOSE 8000

# Run as non-root for security
USER nonroot:nonroot

ENTRYPOINT ["/portfolio"]
