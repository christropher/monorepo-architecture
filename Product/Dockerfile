FROM golang:1.22.2

RUN mkdir /app
ADD . /app
WORKDIR /app

ARG GO_PROXY=https://proxy.golang.org
ENV GOPROXY=$GO_PROXY

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Set the GOOS and GOARCH environment variables
ENV GOOS=linux GOARCH=amd64

RUN go build -o /Product

EXPOSE 8080

ENTRYPOINT [ "/Product" ]