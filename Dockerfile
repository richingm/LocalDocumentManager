FROM golang:1.20.12 as  build-stage

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY go.mod go.sum ./
#RUN go mod download && go mod verify

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mindNotes ./cmd/notes/main.go


FROM alpine:3.19
WORKDIR /app

COPY --from=build-stage /app/mindNotes /app/cmd/notes/mindNotes
COPY --from=build-stage /app/configs  /app/configs
COPY --from=build-stage /app/static /app/static
COPY --from=build-stage /app/templates /app/templates
COPY --from=build-stage /app/data /app/data

EXPOSE 8000

CMD ["/app/cmd/notes/mindNotes"]
