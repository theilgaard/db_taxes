FROM golang:1.22-alpine AS build

RUN apk add build-base

WORKDIR /app
COPY . .

RUN go mod download
RUN go env -w CGO_ENABLED=1 
RUN go build -o /app/db_taxes

FROM alpine:latest

# The image only contains the binary
WORKDIR /app
COPY --from=build /app/db_taxes .

CMD ["./db_taxes"]
