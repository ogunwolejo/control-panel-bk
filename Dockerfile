FROM golang:1.24 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o controlPanelApp ./cmd

RUN chmod +x /app/controlPanelApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/controlPanelApp /app

CMD [ "/app/controlPanelApp" ]


