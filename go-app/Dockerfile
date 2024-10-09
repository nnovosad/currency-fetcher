FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -v -o currency-fetcher .

RUN ls -la

CMD ["/app/currency-fetcher"]
