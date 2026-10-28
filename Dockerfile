FROM golang:1.26

WORKDIR /app

COPY . .

RUN make build

RUN printenv > .env

CMD ["./bin/portfolio-backend"]
