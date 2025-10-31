FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY src/ ./src/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app src/main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/app /app
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app"]