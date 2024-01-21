FROM golang:1.21.6 as build

ENV BIN_FILE /XLabs/app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/*

# Final stage
FROM alpine:3.9

ENV BIN_FILE "/XLabs/app"
ENV CODE_DIR "/XLabs"

WORKDIR ${CODE_DIR}

COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE "/XLabs/configs/config.yaml"
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE}