# STAGE 1: The Builder (The Forge)
FROM golang:1.25.5 AS builder

WORKDIR /app

# 1. Install Build Tools
# Install Templ
RUN go install github.com/a-h/templ/cmd/templ@v0.3.977

# Install Tailwind CSS (Standalone CLI) - Architecture agnostic
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/tag/v4.1.18/download/tailwindcss-linux-x64 && \
    chmod +x tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

# 4. DOWNLOAD TAILWIND v4 (Architecture Aware)
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then \
      URL="https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.18/tailwindcss-linux-arm64"; \
    else \
      URL="https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.18/tailwindcss-linux-x64"; \
    fi && \
    curl -fL -o /usr/local/bin/tailwindcss $URL && \
    chmod +x /usr/local/bin/tailwindcss

# 2. Cache Dependencies
COPY go.mod go.sum ./
RUN ["go", "mod", "download"]

# 3. Copy Source
COPY . .

# 4. GENERATE ASSETS (The Order Matters)
# A. Generate CSS first (so it exists for embedding)
RUN ["tailwindcss", "-i", "./internal/assets/css/input.css", "-o", "./internal/views/css/output.css"]

# B. Generate Templ (which might reference the CSS loader)
RUN ["templ", "generate"]

# 5. Build the Binary
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio ./cmd/web/main.go

# STAGE 2: The Runner (The Vessel)
# Distroless is smaller and safer than Alpine (contains no shell)
FROM gcr.io/distroless/static-debian12

WORKDIR /

# 1. Copy the Binary
COPY --from=builder /app/portfolio /portfolio

# 2. Copy Static Assets
#    We STILL need the assets folder for images and the Godot game files.
COPY --from=builder /app/assets /assets

EXPOSE 8002

# Run as non-root for security
USER nonroot:nonroot

ENTRYPOINT ["/portfolio"]
