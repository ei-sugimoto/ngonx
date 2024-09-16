FROM   golang:1.23.1-alpine AS builder

WORKDIR /app
EXPOSE 8080

RUN --mount=source=go.mod,target=go.mod \
    --mount=source=go.sum,target=go.sum \
    go mod download

RUN --mount=source=.,target=.\
    go build -o /go/bin/main .


FROM alpine:3.12

RUN apk add --no-cache curl
COPY --from=builder /go/bin/main /go/bin/main

CMD ["/go/bin/main"]
