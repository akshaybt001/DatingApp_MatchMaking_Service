FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o match-service

FROM busybox:latest

WORKDIR /match-service/cmd

COPY --from=build /app/cmd/match-service .

COPY --from=build /app/.env /match-service

EXPOSE 8084

CMD ["./match-service"]