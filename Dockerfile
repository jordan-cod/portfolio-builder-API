FROM golang:1.24

WORKDIR /app

COPY . .

RUN make build

CMD ["./bin/portfolio-backend"]
