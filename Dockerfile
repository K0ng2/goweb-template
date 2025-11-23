# Build backend
FROM golang:1.25-alpine AS backend-deps
WORKDIR /app
# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata
# Copy go mod files and download dependencies (leverages Docker layer caching)
COPY backend/go.mod backend/go.sum ./
RUN go mod download && go mod verify

FROM oven/bun:1-alpine AS frontend-deps
WORKDIR /app

COPY web/package.json web/bun.lock ./

# use ignore-scripts to avoid builting node modules like better-sqlite3
RUN bun install --frozen-lockfile --ignore-scripts


FROM oven/bun:1-alpine AS frontend-build
WORKDIR /app

# Copy node_modules from frontend-deps
COPY --from=frontend-deps /app/node_modules ./node_modules
COPY --from=frontend-deps /app/package.json /app/bun.lock ./

# Copy the web source code
COPY web/ ./

RUN bun --bun run generate

FROM backend-deps AS backend-build
# Copy all source code from backend directory
COPY backend/ .
RUN rm -rf backend/web/public
COPY --from=frontend-build /app/.output /app/web

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o app main.go

FROM alpine:latest AS release-base
RUN apk add --no-cache ffmpeg
# Copy timezone data and CA certificates from alpine
COPY --from=backend-build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=backend-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Final stage - minimal runtime image
FROM release-base AS release
# Copy the binary
COPY --from=backend-build /app/app /bin/app

ENV UID=1000
ENV GID=1000
ENV TZ=UTC

# Use non-root user for security
USER $UID:$GID
EXPOSE 8080
ENTRYPOINT ["/bin/app"]
