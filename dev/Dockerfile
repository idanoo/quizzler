FROM golang:1.25-alpine

RUN go install github.com/air-verse/air@latest

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY src/ ./
RUN go mod tidy

EXPOSE 8132

CMD ["air", "-c", ".air.toml"]
