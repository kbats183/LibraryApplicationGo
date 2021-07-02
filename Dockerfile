FROM golang:latest 
ADD . /
WORKDIR /src 
RUN go build -o main .
WORKDIR /
CMD ["/src/main"]

