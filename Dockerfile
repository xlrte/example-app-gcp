FROM golang:1.17.2-alpine3.14 AS builder

RUN apk update && apk add git
WORKDIR /build
COPY . .
RUN ls -ltr
RUN go build cmd/main.go

FROM alpine:3.14

RUN adduser -S hello

RUN apk update
WORKDIR /app

COPY --from=builder /build/main /app/hello

USER hello
CMD ["./hello"]
