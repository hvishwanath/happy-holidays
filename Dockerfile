# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN apk --no-cache add build-base git bzr mercurial gcc
RUN go get -u github.com/gobuffalo/packr/packr
RUN cd /src && packr build

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/happy-holidays /app/
ENTRYPOINT ./happy-holidays
