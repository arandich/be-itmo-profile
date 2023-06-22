FROM golang:alpine AS build

RUN apk --no-cache add gcc g++ make git

WORKDIR /go/src/app

COPY ./app/. .

RUN go get ./...

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./main.go

FROM alpine:3.9

WORKDIR /usr/bin

COPY --from=build /go/src/app/bin /go/bin

EXPOSE 5242

ENTRYPOINT /go/bin/web-app --port 5242