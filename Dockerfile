FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

EXPOSE 7540

ENV TODO_PASSWORD=12345

ENV TODO_PORT=7540

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_todo_app

CMD ["/my_todo_app"]