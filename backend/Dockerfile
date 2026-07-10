# ---- Build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend-api .

# ---- Runtime stage ----
FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /backend-api /backend-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/backend-api"]
