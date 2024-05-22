FROM golang:1.14.10 AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o go-e-commerce_product-api

FROM alpine:latest AS production
COPY --from=builder /app .
ENTRYPOINT ["./go-e-commerce_product-api"]
CMD ["start"]
