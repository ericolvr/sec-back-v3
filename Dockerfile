FROM golang:1.23-alpine AS builder

WORKDIR /src

ENV GOTOOLCHAIN=go1.23.12

COPY go.mod go.sum ./
RUN go version && go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/nr1-back-api ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/nr1-back-api /app/nr1-back-api

ENV PORT=8080

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/nr1-back-api"]
