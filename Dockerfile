FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

WORKDIR /builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /builder/build

FROM alpine:3.21
    
WORKDIR /app
COPY --from=builder --chmod=755 /builder/build /app/app
ENTRYPOINT ["/app/app"]
