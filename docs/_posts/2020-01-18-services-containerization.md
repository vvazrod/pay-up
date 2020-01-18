---
layout: post
title: "Contenerización de los servicios"
---

Definimos imágenes Docker de los distintos servicios que conforman nuestra aplicación para poder desplegarlos de manera más sencilla y en multitud de plataformas. Vamos a ver en este post la definición de la imagen de `gmicro` como ejemplo.

> __Nota:__ Las imágenes definidas para el proyecto se pueden encontrar en [el repositorio de Docker Hub](https://hub.docker.com/r/varrrro/pay-up).

A continuación, vamos a explicar las distintas partes del `Dockerfile` del contenedor.

```dockerfile
FROM golang:1.13-alpine3.10

LABEL maintainer="Víctor Vázquez <victorvazrod@correo.ugr.es>"

WORKDIR /app
```

Usaremos como imagen base la oficial de Go 1.13 hecha con Alpine. En la sección de comparación de imágenes se explica por qué se ha elegido esta en concreto. Definimos quien es el encargado de mantener la imagen y el directorio de trabajo dentro del contenedor, el cuál será `/app`.

```dockerfile
RUN apk update \
    && apk add --no-cache supervisor \
    && apk add --no-cache --virtual .build-deps \
        git \
        curl \
        bash \
    && curl -sL https://git.io/tusk | bash -s -- -b /usr/local/bin latest \
    && apk del .build-deps
```

Existen algunos paquetes que necesitamos instalar en el contenedor para la construcción del mismo. Primero, actualizamos los repositorios de paquetes con `apk update` e instalamos `supervisor`, herramienta que usamos para lanzar y gestionar el proceso del microservicio.

Luego, instalamos `git`, `curl` y `bash` haciendo uso de la opción `--virtual .build-deps`. Con este argumento, lo que hacemos es agrupar estos paquetes instalados dentro de un paquete virtual de nombre `.build-deps` que facilita su gestión. El argumento `--no-cache` hace que no se almacene el índice de paquetes localmente, lo que nos permite reducir el tamaño de la imagen.

Los paquetes de `.build-deps` los necesitamos solo para instalar `tusk`, la herramienta de construcción del proyecto, por lo que los eliminamos después de hacerlo haciendo uso del paquete virtual.

```dockerfile
COPY tusk.yml go.mod go.sum ./

RUN tusk install
```

Copiamos tanto el archivo que define las distintas _tasks_ de `tusk` como los ficheros que especifican las dependencias del proyecto, las cuales se instalan con `tusk install`.

```dockerfile
COPY cmd/gmicro/gmicro.go .
COPY internal/gmicro/*.go internal/gmicro/
COPY internal/gmicro/group/*.go internal/gmicro/group/
COPY internal/gmicro/member/*.go internal/gmicro/member/

ENV CGO_ENABLED=0

RUN tusk build

RUN rm -f gmicro.go && rm -rf internal/
```

Copiamos todos los fuentes del microservicio y compilamos con `tusk build`. Hay que destacar que hay que deshabilitar el uso de CGO para la compilación. CGO es una herramienta que permite la llamada a código C desde un paquete escrito en Go.

Esto ocurre con el paquete `net`, el cuál usamos para crear el servidor HTTP de nuestro servicio. Deshabilitamos entonces el CGO para evitar que el compilador genere un binario dinámico, el cuál provoca errores dentro del contenedor.

Después de compilar el código, eliminamos todos los fuentes para reducir el tamaño de la imagen.

```dockerfile
COPY init/gmicro.conf /etc/supervisor/conf.d/

EXPOSE ${PORT}

CMD [ "tusk", "run" ]
```

Por último, copiamos el fichero de configuración del proceso que usará `supervisord` para lanzarlo y controlarlo, exponemos el puerto del contenedor que va a recibir las peticiones y, con `CMD`, definimos el comando que se ejecuta dentro del contenedor cuando éste se lance (con `docker run`, por ejemplo). En nuestro caso, `tusk run` iniciará el microservicio a través de `supervisord`.

## Actualización a construcción multifase

Al ser Go un lenguaje compilado, podemos aplicar una estrategia de construcción multifase para reducir el tamaño de la imagen al mínimo. La idea es compilar la aplicación en un contenedor y copiar el binario obtenido a otro contenedor, que contendrá solo este binario.

Vamos a echar un vistazo al `Dockerfile` actualizado.

```dockerfile
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
```

Esta es la fase de construcción, definida con la etiqueta _build_ en la primera línea. Como se puede apreciar, el proceso es prácticamente idéntico al que se realizaba en el `Dockerfile` anterior, con la excepción de que ya no nos preocupamos de eliminar los archivos y demás utilidades instaladas en la imagen ya que no es necesario. Tampoco definimos un `CMD` o `ENTRYPOINT`, sino que la imagen termina compilando la aplicación con `tusk build gmicro`.

```dockerfile
FROM alpine:3.10
LABEL maintainer="Víctor Vázquez <victorvazrod@correo.ugr.es>"
WORKDIR /app

# Copy binary from build stage
COPY --from=build /usr/local/bin/gmicro /app/
ENTRYPOINT ./gmicro
```

Esta sería la imagen "real" que se utiliza para lanzar la aplicación contenerizada. Se usa de imagen base de la de `alpine`, sin necesidad de usar la de `golang` como en la fase de construcción ya que no necesitamos hacer uso del _toolchain_ de Go.

Lo único que se hace en esta imagen es copiar el binario de la imagen de construcción con `COPY --from=build` y definir dicho binario como `ENTRYPOINT` para que se ejecute al lanzar el contenedor.

Usando esta técnica conseguimos reducir el tamaño de la imagen de `gmicro` a 8.53MB, una tamaño sumamente inferior al que teníamos antes.
