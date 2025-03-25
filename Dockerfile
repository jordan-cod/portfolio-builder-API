FROM golang:1.24

WORKDIR /app

COPY . .

RUN make build

RUN printenv > .env

CMD ["./bin/portfolio-backend"]
