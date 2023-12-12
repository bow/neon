# Dockerfile for packaging releases.
#
# Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
# SPDX-License-Identifier: BSD-3-Clause
#
# This file is part of lens <https://github.com/bow/lens>.

FROM golang:1.21-alpine AS builder

WORKDIR /src

RUN apk add --update --no-cache build-base~=0 make~=4 git~=2

COPY .git /src/.git

RUN git checkout -- . && make bin

# -- #

FROM golang:1.21-alpine

ARG REVISION
ARG BUILD_TIME

LABEL org.opencontainers.image.title="lens" \
    org.opencontainers.image.url="https://ghcr.io/bow/lens" \
    org.opencontainers.image.source="https://github.com/bow/lens" \
    org.opencontainers.image.authors="Wibowo Arindrarto <contact@arindrarto.dev>" \
    org.opencontainers.image.revision="${REVISION}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.licenses="BSD-3-Clause"

COPY --from=builder /src/bin/lens /bin/lens

RUN mkdir -p /var/data/
ENV LENS_SERVE_ADDR=tcp://0.0.0.0:7000 \
    LENS_SERVE_DB_PATH=/var/data/lens.db

WORKDIR /runtime
ENTRYPOINT ["lens"]
