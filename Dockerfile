FROM golang
ADD . /go/src/github.com/tobyjsullivan/event-history
RUN  go install github.com/tobyjsullivan/event-history

EXPOSE 3000

CMD /go/bin/event-history
