FROM golang:1.26-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o issue-tracker ./cmd/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/issue-tracker .

EXPOSE 3002

CMD [ "./issue-tracker" ]