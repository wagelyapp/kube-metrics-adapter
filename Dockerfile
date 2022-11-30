FROM golang:1.19 AS BUILDER

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
ADD . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o kube-metrics-adapter .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=BUILDER /app/kube-metrics-adapter /

ENTRYPOINT ["/kube-metrics-adapter"]
