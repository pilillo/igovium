ARG ARTIFACT=igovium
FROM golang:1.16-alpine AS builder

ARG ARTIFACT

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ${ARTIFACT} .

# ---->

FROM alpine:3.12.4

ARG ARTIFACT
ENV ARTIFACT=${ARTIFACT}

# olric
# EXPOSE 3320 3322

# set default vars
ENV GIN_MODE=release
# do not forget to set the config and copy the conf file
# ENV IGOVIUM_CONFIG=/conf.yaml
# COPY conf.yaml /
COPY --from=builder /build/${ARTIFACT} ./

# Command to run when starting the container
ENTRYPOINT ["sh", "-c", "./${ARTIFACT}"]