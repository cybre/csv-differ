FROM golang:1.21.4-alpine3.18 as builder
WORKDIR /app
COPY ./ ./
RUN go build ./cmd/csv-differ

FROM alpine:3.18 AS prod
WORKDIR /app
COPY --from=builder /app/csv-differ /app/
ENTRYPOINT ["/app/csv-differ"]
