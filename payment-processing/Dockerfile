### STAGE 1 : Build the go source code into binary
FROM golang:1.17 as builder

ENV APP_DIR /payment-processing
ENV GOFLAGS -mod=vendor

## Copy source code from local machine into container
RUN mkdir -p ${APP_DIR}
COPY . ${APP_DIR}

RUN cd $APP_DIR && go test $(go list ./... | grep -v /vendor/)

# Compile the binary and statically link
RUN cd $APP_DIR && CGO_ENABLED=0 go build -ldflags '-d -w -s'

### STAGE 2 : Package the binary in a minimal alpine base image
FROM alpine:latest

ENV APP_DIR /payment-processing

COPY --from=builder ${APP_DIR}/payment-processing .
RUN mkdir config
COPY --from=builder ${APP_DIR}/config/config.env ./config

RUN apk add tzdata ca-certificates

EXPOSE 8080

CMD ["./payment-processing"]