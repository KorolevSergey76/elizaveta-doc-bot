FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/bot

FROM alpine:3.21
WORKDIR /app
COPY --from=build /app /app/bot
CMD ["/app/bot"]
