FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN go build -o /bin/prediction-api ./cmd/api

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /bin/prediction-api /app/prediction-api
COPY --from=build /app/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/prediction-api"]
