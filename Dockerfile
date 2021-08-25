# ------------------------------------------------------------------------------
# Build image
# ------------------------------------------------------------------------------
FROM golang:1.16-alpine3.13 as build

# install build tools
RUN set -eux; \
	apk add -U --no-cache \
		curl \
		git  \
		make \
	;

# add new user
ARG USER=default
RUN set -eux; \
    addgroup ${USER}; \
    adduser -h /build -D -G ${USER} ${USER};

USER ${USER}
WORKDIR /build

# copy vendored dependencies
COPY --chown=${USER}:${USER} go.mod go.sum ./
COPY --chown=${USER}:${USER} ./vendor ./vendor
COPY --chown=${USER}:${USER} ./internal ./internal
COPY --chown=${USER}:${USER} ./cmd ./cmd
COPY --chown=${USER}:${USER} ./Makefile ./Makefile

RUN set -eux; \
	make build;

# ------------------------------------------------------------------------------
# Base runtime image
# ------------------------------------------------------------------------------
FROM alpine:3.13 AS base-runtime

# add new user
ARG USER=default
RUN set -eux; \
    addgroup ${USER}; \
    adduser -D -G ${USER} ${USER};

USER ${USER}
WORKDIR /home/${USER}

# ------------------------------------------------------------------------------
# pxe-init-server Runtime image
# ------------------------------------------------------------------------------
FROM base-runtime as pxe-init-server

COPY --from=build /build/dist/pxe-init-server /usr/bin/pxe-init-server

CMD ["/usr/bin/pxe-init-server", "serve"]