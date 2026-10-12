FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /scheduler .

FROM ubuntu:latest
WORKDIR /app
COPY --from=build /scheduler ./scheduler
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db

EXPOSE 7540

CMD ["./scheduler"]
