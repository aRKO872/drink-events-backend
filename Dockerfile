FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /drink-events-go

EXPOSE 3050
CMD air