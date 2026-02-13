# STAGE 1: The Builder (The Forge)
FROM golang:1.25.5 AS builder

WORKDIR /app

# 1. Install Build Tools
# Install Templ
RUN go install github.com/a-h/templ/cmd/templ@v0.3.977

# Install Tailwind CSS (Standalone CLI)
# Download the linux binary so we can compile CSS without Node.js
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

# 2. Cache Dependencies
COPY go.mod go.sum ./
RUN ["go", "mod", "download"]

# 3. Copy Source & Generate
COPY . .

# Generate the HTML templates
RUN ["templ", "generate"]

# Generate the CSS 
# Reads from internal/ and outputs to assets/
RUN ["tailwindcss", "-i", "./internal/assets/css/input.css", "-o", "./assets/css/output.css"]

# 4. Test & Build
RUN ["go", "test", "-v",  "./..."] 
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio ./cmd/web/main.go

# STAGE 2: The Runner (The Vessel)
FROM gcr.io/distroless/static-debian12

WORKDIR /

# 1. Copy the Binary
COPY --from=builder /app/portfolio /portfolio

# 2. Copy the Assets 
COPY --from=builder /app/assets /assets

EXPOSE 8000

# Run as non-root for security
USER nonroot:nonroot

ENTRYPOINT ["/portfolio"]
