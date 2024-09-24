FROM golang:1.22.6-alpine

WORKDIR /app
COPY . .
ENTRYPOINT ["go", "run", "."]
