############################
# STEP 1 build executable binary
############################
FROM golang:1.21-alpine AS builder

WORKDIR /go/src/app
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,target=. \
    go install

############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
WORKDIR /root/
RUN apk add --no-cache ffmpeg
COPY --from=builder /go/bin/madcap .
COPY .env.vault .
CMD ["./madcap"]
