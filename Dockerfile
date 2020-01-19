#--------------------#
# Build stage
#--------------------#
FROM golang:1.13-alpine3.10 AS build

# Install needed utilities
RUN apk update \
    && apk add --no-cache supervisor git curl bash \
    && curl -sL https://git.io/tusk | bash -s -- -b /usr/local/bin latest

# Copy task runner config and module files
COPY tusk.yml go.mod go.sum /src/

# Install project dependencies
RUN cd /src && tusk install

# Copy source files
COPY cmd/gmicro/main.go /src/cmd/gmicro/
COPY internal/gmicro/ /src/internal/gmicro
COPY internal/consumer/ /src/internal/consumer/
COPY internal/tmicro/expense/ /src/internal/tmicro/expense/
COPY internal/tmicro/payment/ /src/internal/tmicro/payment/

# Disable CGO
ENV CGO_ENABLED=0

# Build binary
RUN cd /src && tusk build gmicro

#--------------------#
# Deployment stage
#--------------------#
FROM alpine:3.10
LABEL maintainer="Víctor Vázquez <victorvazrod@correo.ugr.es>"
WORKDIR /app

# Copy binary from build stage
COPY --from=build /usr/local/bin/gmicro /app/
ENTRYPOINT ./gmicro