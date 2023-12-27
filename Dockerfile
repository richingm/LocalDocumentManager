FROM golang:1.20.12

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY go.mod go.sum ./
#RUN go mod download && go mod verify

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -v -o /app/local_document_manager ./cmd/local_document_manager/main.go

CMD ["/app/local_document_manager"]
