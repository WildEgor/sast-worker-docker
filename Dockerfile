# Base Stage
FROM golang:1.22-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache  \
    wget  \
    bash \
    tar \
    ca-certificates \
    && update-ca-certificates
# prepare hadolint
RUN wget -O /usr/local/bin/hadolint  https://github.com/hadolint/hadolint/releases/latest/download/hadolint-Linux-x86_64 \
    && chmod +x /usr/local/bin/hadolint
# prepare trivy
RUN wget https://github.com/aquasecurity/trivy/releases/latest/download/trivy_0.52.2_Linux-64bit.tar.gz \
    && tar zxvf trivy_0.52.2_Linux-64bit.tar.gz \
    && mv trivy /usr/local/bin/ \
    && rm trivy_0.52.2_Linux-64bit.tar.gz
# RUN hadolint --version
# RUN trivy --version
RUN mkdir -p temp \
    mkdir -p dist
RUN go mod download

# Development Stage
FROM base as dev
WORKDIR /app/
COPY . .
# for hot-reloading
RUN go install github.com/air-verse/air@latest
RUN go build -o dist/app cmd/main.go
CMD ["air", "-c", ".air-unix.toml", "-d"]

# Debug stage
FROM base as debug
WORKDIR /
COPY . .
RUN mkdir -p /app/temp
COPY ./scripts /app/scripts
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go build -gcflags="all=-N -l" -o /app/app cmd/main.go
RUN mv /go/bin/dlv /
EXPOSE 8081 40000
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/app"]

# Build Production Stage
FROM base as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o dist/app cmd/main.go

# Production Stage
FROM cgr.dev/chainguard/bash:latest as production
WORKDIR /app/
COPY --from=build /app/temp /app/temp
COPY --from=build /app/scripts /app/scripts
COPY --from=build /app/dist/app .
# Specify method fetch config.yaml!
COPY --from=base /usr/local/bin/hadolint /usr/local/bin/hadolint
COPY --from=base /usr/local/bin/trivy /usr/local/bin/trivy
CMD ["./app"]