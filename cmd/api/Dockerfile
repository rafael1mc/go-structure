FROM golang:1.21.5-alpine as build

WORKDIR /app-build

COPY go.mod .
COPY go.sum .

RUN go mod download

# TODO will allow for auto reload inside container
# COPY .air.toml .

COPY keys keys
# COPY seed seed
# COPY migrations migrations

COPY internal internal
COPY cmd/api cmd/api

RUN go build -o app ./cmd/api

FROM alpine

WORKDIR /app

# COPY --from=build /app-build/.air.toml .
COPY --from=build /app-build/keys keys
# COPY --from=build /app-build/seed seed
# COPY --from=build /app-build/migrations migrations
COPY --from=build /app-build/app .

CMD ["./app"]