FROM golang:1.15-alpine as build

WORKDIR /src

RUN apk add -U --no-cache ca-certificates && \
    apk add build-base git fuse fuse-dev

ARG ACCESS_TOKEN_USR="mallelaprasanth"
ARG ACCESS_TOKEN_PWD="2ab0ddecaa5d77ff50de8c8468540be7f5bc0453"

RUN printf "machine github.com\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n\
    \n\
    machine api.github.com\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n"\
    >> /root/.netrc
RUN chmod 600 /root/.netrc

# COPY .gitconfig /root/.gitconfig
COPY go.mod ./
COPY go.sum ./
COPY cmd ./
COPY entrypoint.sh ./
COPY public/index.html ./public/index.html

ENV GOPRIVATE="github.com/synspective"


RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/mvt-server ./main.go

# Install GCFuse
WORKDIR /gcsfuse

ARG GCSFUSE_VERSION=0.30.0

RUN go get -d github.com/googlecloudplatform/gcsfuse
RUN go install github.com/googlecloudplatform/gcsfuse/tools/build_gcsfuse

RUN build_gcsfuse ${GOPATH}/src/github.com/googlecloudplatform/gcsfuse /tmp ${GCSFUSE_VERSION}

# running image
FROM alpine

WORKDIR /go

RUN apk --update add fuse \
    && rm -rf /var/cache/apk/*

ARG GCSFUSE_VERSION=0.30.0

ARG PROJECT_ID
ENV PROJECT_ID $PROJECT_ID

COPY --from=build /src/entrypoint.sh /go/
COPY --from=build /src/bin /go/bin
COPY --from=build /src/public/index.html /go/public/index.html

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /tmp/bin/gcsfuse /usr/bin
COPY --from=build /tmp/sbin/mount.gcsfuse /usr/sbin

ENTRYPOINT ["/go/entrypoint.sh"]
