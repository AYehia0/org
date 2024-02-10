FROM golang:1.21.5-alpine3.17 AS build
WORKDIR /app

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:3.17
WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build /app/config /app/config
COPY --from=build /app/*.sh /app/
COPY --from=build /app/fullchain.pem /app/
COPY --from=build /app/privkey.pem /app/

CMD ["/app/main"]
ENTRYPOINT ["/app/run.sh"]
