FROM golang:1.12.7-stretch AS builder

ENV GOBIN /go/bin
ENV GO111MODULE on 
ENV CGO_ENABLED 0

RUN mkdir /app/
ADD . /app/
WORKDIR /app/

RUN go build -o bin/main main.go

FROM alpine
WORKDIR /app
RUN mkdir /app/app
COPY --from=builder /app/bin/main /app/main

CMD ["/app/main"]
