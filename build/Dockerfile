FROM golang:1.14.2 as builder

WORKDIR /go/src/github.com/mrbttf/bridge-server 

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go install -v /go/src/github.com/mrbttf/bridge-server/cmd/main.go


FROM scratch

COPY --from=builder /go/bin/main /go/bin/main

ENTRYPOINT ["/go/bin/main"]

EXPOSE 8080