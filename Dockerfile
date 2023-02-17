# Dockerfile for packaging releases.
#
# Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
# SPDX-License-Identifier: BSD-3-Clause
#
# This file is part of Iris <https://github.com/bow/iris>.

FROM golang:1.20-alpine AS builder

WORKDIR /src

RUN apk add --update --no-cache build-base~=0 make~=4 git~=2

COPY .git /src/.git

RUN git checkout -- . && make bin

# -- #

FROM golang:1.20-alpine

ARG REVISION
ARG BUILD_TIME

LABEL org.opencontainers.image.title="iris" \
    org.opencontainers.image.url="https://ghcr.io/bow/iris" \
    org.opencontainers.image.source="https://github.com/bow/iris" \
    org.opencontainers.image.authors="Wibowo Arindrarto <contact@arindrarto.dev>" \
    org.opencontainers.image.revision="${REVISION}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.licenses="BSD-3-Clause"

COPY --from=builder /src/bin/iris /bin/iris

RUN mkdir -p /var/data/
ENV IRIS_IN_DOCKER=1 \
    IRIS_SERVE_ADDR=tcp://0.0.0.0:7000 \
    IRIS_SERVE_DB_PATH=/var/data/iris.db

WORKDIR /runtime
ENTRYPOINT ["iris"]
