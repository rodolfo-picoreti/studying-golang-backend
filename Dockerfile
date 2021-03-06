FROM golang:1.17-alpine3.15 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN GOOS=linux CGO_ENABLED=0 go build -o app

###

FROM alpine:3.15

COPY --from=builder /app/app /app

CMD ["/app"]

