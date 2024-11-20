FROM golang:1.23-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o /transactions-manager-app

FROM gcr.io/distroless/base-debian11 AS runtime

COPY --from=builder /transactions-manager-app /transactions-manager-app

ENV PORT=8000

EXPOSE $PORT

CMD ["/transactions-manager-app"]
