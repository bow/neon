FROM golang:1.18-alpine AS builder

WORKDIR /src

RUN apk add --update --no-cache build-base~=0 make~=4 git~=2

COPY .git /src/.git

RUN git checkout -- . && make bin

# -- #

FROM golang:1.18-alpine

ARG REVISION
ARG BUILD_TIME

LABEL org.opencontainers.image.title="courier" \
    org.opencontainers.image.url="https://ghcr.io/bow/courier" \
    org.opencontainers.image.source="https://github.com/bow/courier" \
    org.opencontainers.image.authors="Wibowo Arindrarto <contact@arindrarto.dev>" \
    org.opencontainers.image.revision="${REVISION}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.licenses="BSD-3-Clause"

COPY --from=builder /src/bin/courier /bin/courier

RUN mkdir -p /var/data/
ENV COURIER_IN_DOCKER=1 \
    COURIER_SERVE_ADDR=tcp://0.0.0.0:7000 \
    COURIER_SERVE_DB_PATH=/var/data/courier.db

WORKDIR /runtime
ENTRYPOINT ["courier"]
