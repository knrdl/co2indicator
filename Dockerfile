FROM golang:alpine as builder

WORKDIR /usr/src/app
COPY src .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /co2indicator



FROM scratch

EXPOSE 8080/tcp

COPY --from=builder /co2indicator /co2indicator

ENTRYPOINT ["/co2indicator", "--server", ":8080"]
