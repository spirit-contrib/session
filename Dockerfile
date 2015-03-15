FROM google/golang

WORKDIR /gopath/src/github.com/spirit-contrib/session
ADD . /gopath/src/github.com/spirit-contrib/session/
RUN go get github.com/spirit-contrib/session

CMD []
ENTRYPOINT []