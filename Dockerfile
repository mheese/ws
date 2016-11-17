FROM golang:1.7.3

ADD . /go/src/github.com/mheese/ws/

CMD [ "/go/src/github.com/mheese/ws/build.sh" ]
