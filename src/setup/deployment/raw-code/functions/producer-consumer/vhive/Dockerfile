FROM docker.io/vhiveease/vhive-golang:latest as build

WORKDIR /app

COPY go.mod /app/go.mod
COPY main.go /app/main.go
COPY common /app/common/
COPY proto_gen /app/proto_gen/

RUN go mod download && \
   CGO_ENABLED=0 GOOS=linux go build -v -o /main /app/main.go

FROM scratch as bin-unix
COPY --from=build /main /main
COPY filler.file /filler.file
ENTRYPOINT [ "/main" ]