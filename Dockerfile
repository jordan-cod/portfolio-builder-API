FROM golang:1.21

WORKDIR /app

COPY . .

RUN make build

CMD ["$(BIN_DIR)/$(APP_NAME)"]
