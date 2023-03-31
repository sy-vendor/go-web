FROM golang:1.18 AS builder

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o apiserver ./cmd/apiserver

FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /src/apiserver /app/

ENTRYPOINT ["./apiserver"]