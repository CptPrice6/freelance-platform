FROM golang:1.24 as builder

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o backend

RUN ls -l /app/backend

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/backend /root/

COPY conf ./conf

EXPOSE 8080

# Run the application
CMD ["/root/backend"]
