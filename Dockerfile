# Dockerfile for packaging releases.
#
# Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
# SPDX-License-Identifier: BSD-3-Clause
#
# This file is part of neon <https://github.com/bow/neon>.

FROM golang:1.21-alpine AS builder

WORKDIR /src

RUN apk add --update --no-cache build-base~=0 make~=4 git~=2

COPY .git /src/.git

RUN git checkout -- . && make bin

# -- #

FROM golang:1.21-alpine

ARG REVISION
ARG BUILD_TIME

LABEL org.opencontainers.image.title="neon" \
    org.opencontainers.image.url="https://ghcr.io/bow/neon" \
    org.opencontainers.image.source="https://github.com/bow/neon" \
    org.opencontainers.image.authors="Wibowo Arindrarto <contact@arindrarto.dev>" \
    org.opencontainers.image.revision="${REVISION}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.licenses="BSD-3-Clause"

COPY --from=builder /src/bin/neon /bin/neon

RUN mkdir -p /var/data/
ENV NEON_SERVE_ADDR=tcp://0.0.0.0:5151 \
    NEON_SERVE_DB_PATH=/var/data/neon.db

WORKDIR /runtime
ENTRYPOINT ["neon"]
