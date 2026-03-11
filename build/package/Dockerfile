FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN "CGO_ENABLED=0 GOOS=linux go build -o /bin/migr8"

EXPOSE 8080

ENTRYPOINT [ "/bin/migr8" ]
