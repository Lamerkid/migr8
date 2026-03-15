FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/migr8 ./cmd/migr8

ENV PATH "/app/bin:${PATH}"

ENTRYPOINT [ "migr8" ]
