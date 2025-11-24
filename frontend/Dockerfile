FROM golang:1.25-bookworm AS builder

WORKDIR /build/

COPY . /build/

RUN go build -tags netgo -ldflags '-s -w' -o server

FROM debian:bookworm-slim
WORKDIR /app/
COPY --from=builder /build/server /app/
RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
COPY --from=builder /build/static/ /app/static/

EXPOSE 8000

ENTRYPOINT [ "./server" ]