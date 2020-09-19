FROM golang:1.15 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN chmod +x /app/wait-for-it.sh
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/
RUN CGO_ENABLED=0 GOOS=linux go build -o admin ./cmd/admin/

FROM alpine:3.10 AS production
RUN apk --no-cache add ca-certificates --upgrade bash


COPY --from=builder /app .

CMD ["sh","run.sh"]


