FROM golang:1.14 AS GO_BUILD
ENV CGO_ENABLED 0
COPY . /lemon_cash
WORKDIR /lemon_cash
RUN go mod download
RUN go build -o server

FROM alpine:3.15
WORKDIR /lemon_cash
COPY --from=GO_BUILD /lemon_cash/server /lemon_cash
EXPOSE 8080
CMD ["./server"]