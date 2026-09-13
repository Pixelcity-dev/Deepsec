FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /deepsec ./cmd/deepsec

FROM gcr.io/distroless/static-debian12:latest

COPY --from=builder /deepsec /usr/local/bin/deepsec

COPY rules /rules

ENTRYPOINT ["deepsec"]
