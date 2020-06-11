FROM golang:1.14.2

WORKDIR /go/src/github.com/mrbttf/bridge-server 

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go install -v /go/src/github.com/mrbttf/bridge-server/cmd        

EXPOSE 8080

CMD [ "/go/bin/cmd" ]