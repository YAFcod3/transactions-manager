FROM golang:1.23-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG ENV=production
ENV ENV=$ENV

RUN if [ "$ENV" = "development" ]; then go test -v ./...; fi

RUN CGO_ENABLED=0 GOOS=linux go build -o /transactions-manager-app

FROM gcr.io/distroless/base-debian11 AS runtime

COPY --from=builder /transactions-manager-app /transactions-manager-app

ENV ENV=production
ENV PORT=8000

EXPOSE $PORT

CMD ["/transactions-manager-app"]
