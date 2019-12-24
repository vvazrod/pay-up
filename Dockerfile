FROM golang:1.13-alpine3.10

LABEL maintainer="Víctor Vázquez <victorvazrod@correo.ugr.es>"

WORKDIR /app

# Install supervisor and other build dependencies, clean up after the fact
RUN apk update \
    && apk add --no-cache supervisor \
    && apk add --no-cache --virtual .build-deps \
        git \
        curl \
        bash \
    && curl -sL https://git.io/tusk | bash -s -- -b /usr/local/bin latest \
    && apk del .build-deps

# Copy task runner config and module files
COPY tusk.yml go.mod go.sum ./

# Install project dependencies
RUN tusk install

# Copy source files
COPY cmd/gmicro/gmicro.go .
COPY internal/gmicro/*.go internal/gmicro/
COPY internal/gmicro/group/*.go internal/gmicro/group/
COPY internal/gmicro/member/*.go internal/gmicro/member/

# Disable CGO
ENV CGO_ENABLED=0

# Compile source files
RUN tusk build

# Copy service supervisor config
COPY init/gmicro.conf /etc/supervisor/conf.d/

EXPOSE ${PORT}

# Add non-root user used to run supervisord process
RUN addgroup -S gmicro \
    && adduser -S -G gmicro gmicro \
    && chown gmicro:gmicro /app

CMD [ "tusk", "run" ]