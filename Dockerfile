FROM golang:alpine

ENV \
    HTTP_LISTEN=":4080" \
    PACKAGE_PATH="github.com" \
    CLONE_PATH="https://github.com"

WORKDIR /app

COPY ./ /app

RUN go build -o meta-go-imports main.go

CMD [ "./meta-go-imports" ]
