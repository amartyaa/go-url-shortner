FROM golang:alpine AS builder

RUN  mkdir /builder

ADD . /build/

WORKDIR /build

RUN  go build -o main .

#Stage 2 
FROM alpine:latest

RUN adduser -D -g '' appuser

USER appuser

COPY . /app

COPY --from=builder /build/main .

WORKDIR /app

EXPOSE 8080

CMD ["./main"]