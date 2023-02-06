FROM golang

WORKDIR /app

COPY . ./

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o main main.go

CMD [ "./main" ]
