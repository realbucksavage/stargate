FROM golang:1.19-alpine AS build

WORKDIR /go/src/github.com/realbucksavage/stargate

COPY . .
RUN go build -o ./app _examples/websockets/main.go

FROM alpine
COPY --from=build /go/src/github.com/realbucksavage/stargate/app ./app
ENTRYPOINT ["./app"]