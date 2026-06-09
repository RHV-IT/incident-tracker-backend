FROM golang:1.22-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o issue-tracker ./cmd/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/issue-tracker .
COPY --from=builder /app/tables.sql .
COPY --from=builder /app/createtables.sh .

RUN sh createtables.sh

EXPOSE 3002

CMD [ "./issue-tracker" ]