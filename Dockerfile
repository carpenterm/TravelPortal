FROM golang:1.25-alpine AS builder

ARG GOPROXY="https://proxy.golang.org"
ENV GOPROXY=${GOPROXY}
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
# cache modules during CI-friendly builds
RUN go env -w GOPROXY=${GOPROXY} && go mod download

COPY . .

# Install C build dependencies required by mattn/go-sqlite3
RUN apk add --no-cache gcc musl-dev linux-headers sqlite-dev

# Build with CGO enabled for sqlite3 support
ENV CGO_ENABLED=1
RUN go build -ldflags "-s -w" -o /app/travelportal .

FROM alpine:3.18

LABEL org.opencontainers.image.source="https://github.com/skandpurohit/TravelPortal"

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/travelportal .
COPY --from=builder /app/web ./web

# Create non-root user for safer runtime in CI/CD environments
# A numeric UID, so Kubernetes can verify runAsNonRoot.
RUN addgroup -S -g 10001 app && adduser -S -u 10001 -G app app \
	&& chown -R 10001:10001 /app \
	&& chmod +x /app/travelportal

# Use the fallback variable for the demo default so either password variable can override it.
ENV TODO_PASSWORD=test12345

EXPOSE 7540

USER 10001:10001

CMD ["/app/travelportal"]
