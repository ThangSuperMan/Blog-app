FROM golang:1.18.10-alpine3.17

WORKDIR /app

COPY . ./

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o main main.go

CMD [ "./main" ]
