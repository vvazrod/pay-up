FROM golang:1.13-alpine3.10

# Set the maintainer of this image
LABEL maintainer="Víctor Vázquez <victorvazrod@correo.ugr.es>"

# Establish the working directory
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
COPY cmd/gateway/main.go .
COPY internal/gateway/ internal/gateway
COPY internal/publisher/ internal/publisher/
COPY internal/tmicro/expense/ internal/tmicro/expense/
COPY internal/tmicro/payment/ internal/tmicro/payment/

# Disable CGO
ENV CGO_ENABLED=0

# Compile source files and delete them after the fact
RUN tusk build gateway \
    && rm -f main.go && rm -rf internal/

# Copy service supervisor config
COPY init/gateway.conf /etc/supervisor/conf.d/

# Expose container port 8080
EXPOSE 8080

# Run the application
CMD [ "tusk", "run", "gateway" ]