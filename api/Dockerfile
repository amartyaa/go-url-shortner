FROM golang:alpine AS builder

RUN  mkdir /build

ADD . /build/

WORKDIR /build

RUN  go build -o main .

#Stage 2 
FROM alpine:latest

RUN adduser -D -S -H -h /app appuser

USER appuser

COPY . /app

COPY --from=builder /build/main /app/

WORKDIR /app

EXPOSE 8080

CMD ["./main"] --v