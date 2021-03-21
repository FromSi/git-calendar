FROM golang:alpine as builder

RUN apk --no-cache add ca-certificates git
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine
COPY --from=builder /app/calendar /app/
EXPOSE 2004 
ENTRYPOINT ["/app/calendar"]
