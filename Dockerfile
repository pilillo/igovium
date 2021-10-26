ARG ARTIFACT=igovium
FROM golang:1.16-alpine AS builder

ARG ARTIFACT
ARG MAIN_PATH=main.go

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ${ARTIFACT} ${MAIN_PATH}

# ---->
FROM alpine:3.12.4

ARG ARTIFACT
ENV ARTIFACT=${ARTIFACT}

# olric
# https://github.com/buraksezer/olric/blob/master/Dockerfile
EXPOSE 3320 3322

# set default vars
ENV GIN_MODE=release

COPY --from=builder /build/${ARTIFACT} ./

# Command to run when starting the container
ENTRYPOINT ["sh", "-c", "./${ARTIFACT}"]