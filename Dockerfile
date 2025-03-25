FROM golang:1.24

WORKDIR /app

COPY . .

RUN make build

CMD ["$(BIN_DIR)/$(APP_NAME)"]
