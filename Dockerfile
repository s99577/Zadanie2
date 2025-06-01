# Etap 1: Budowanie aplikacji
FROM golang:1.24-alpine@sha256:b4f875e650466fa0fe62c6fd3f02517a392123eea85f1d7e69d85f780e4db1c1 AS builder

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /usr/app

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./

RUN CGO_ENABLED=0 go build -o /usr/app/weather_app ./main.go

FROM scratch AS prod

ENV AUTHOR="Klaudia Klimont"
LABEL org.opencontainers.image.authors=${AUTHOR}

COPY --from=builder /usr/app/ /usr/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE ${PORT}

WORKDIR /usr/app

ENTRYPOINT ["./weather_app"]