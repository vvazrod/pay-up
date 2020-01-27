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
COPY cmd/gateway/main.go /src/cmd/gateway/
COPY internal/gateway/ /src/internal/gateway/
COPY internal/publisher/ /src/internal/publisher/
COPY internal/tmicro/expense/ /src/internal/tmicro/expense/
COPY internal/tmicro/payment/ /src/internal/tmicro/payment/

# Disable CGO
ENV CGO_ENABLED=0

# Build binary
RUN cd /src && tusk build gateway

#--------------------#
# Deployment stage
#--------------------#
FROM alpine:3.10
LABEL maintainer="Víctor Vázquez <victorvazrod@correo.ugr.es>"
WORKDIR /app

# Copy binary from build stage
COPY --from=build /src/gateway /app/
ENTRYPOINT ./gateway