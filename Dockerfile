FROM golang:1.17-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o app

###

FROM alpine

COPY --from=builder /app/app /app

CMD ["/app"]

