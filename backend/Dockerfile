FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o backend

RUN ls -l /app/backend

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/backend /root/

COPY conf ./conf

COPY --from=builder /app/migrations /root/migrations

EXPOSE 8080

# Run the application
CMD ["/root/backend"]
