FROM golang:1.16-alpine

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -d -v
RUN go install -v

RUN go build

# Build commands
RUN mkdir -p build
RUN go build -v -mod=readonly -ldflags="-s -w" -o build ./cmd/...

CMD ["dpe-insights"]
