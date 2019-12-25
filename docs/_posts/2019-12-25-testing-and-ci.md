---
layout: post
title: "Pruebas e integración continua"
---

Al utilizar Go para el proyecto, nos encontramos con la circunstancia de que no hay marcos de prueba para este lenguaje. El propio lenguaje incorpora un *toolchain* con todas las herramientas necesarias para la gestión de un proyecto. Así pues, usaremos el paquete `testing` para escribir las pruebas y la herramienta `go test` para ejecutarlas.

Como *task runner*, usaremos [Tusk](https://github.com/rliebz/tusk), el cual nos permite definir las tareas a ejecutar mediante un fichero `tusk.yml`.

Tradicionalmente, para proyectos escritos en Go se ha utilizado `make`, que es muy versátil pero también muy difícil de leer y escribir. De entre las alternativas más modernas ([mage](https://github.com/magefile/mage) o [Taskfile](https://github.com/go-task/task)) se ha elegido Tusk por ser más sencilla su instalación y uso.

## Configuración de Travis

Vamos a utilizar Travis para la integración continua de nuestro proyecto y, a continuación, vamos a explicar las distintas partes del fichero `.travis.yml` que define nuestras *builds*.

```yaml
language: go

go:
  - 1.11.x
  - 1.13.x
  - master
```

Para empezar, con la etiqueta `language` indicamos que nuestro proyecto está implementado en Go y en la etiqueta `go` indicamos las versiones. Travis lanzará una tarea o *job* para cada una de estas versiones de Go, llevando a cabo una *build* para cada una de ellas.

La versión `master` representa la versión más reciente de Go, mientras que usamos la 1.11 como versión mínima por ser la primera que implementa la funcionalidad de módulos, la cual utilizamos en nuestro proyecto para gestionar las dependencias.

```yaml
before_install:
  - curl -sL https://git.io/tusk | sudo bash -s -- -b /usr/local/bin latest
```

Como su propio nombre indica `before_install` especifica los comandos que se deben ejecutar antes de la fase `install`. En este caso, lo que hacemos es instalar la herramienta Tusk, que necesitaremos para ejecutar las siguientes tareas de la *build*.

```yaml
install:
  - tusk install

script:
  - tusk test

after_success:
  - tusk coverage
```

Estas tareas de Tusk están definidas en el fichero `tusk.yml`.

```yaml
tasks:
  install:
    usage: Install project dependencies
    run: go mod download
  
  test:
    usage: Run tests and generate coverage report
    run: go test ./internal/... -coverprofile coverage.txt

  coverage:
    usage: Send coverage report to codecov.io
    run: bash <(curl -s https://codecov.io/bash)
```

La fase `install` se usa para instalar las dependencias del proyecto. Como podemos ver, lo que se hace con `tusk install` es ejecutar `go mod download`, lo que descarga todas las dependencias deifinas en nuestro fichero `go.mod`.

En la fase `script` es donde se realiza el proceso de *build* como tal de nuestro proyecto. A estas alturas del proyecto, esto consiste solo en la ejecución de los tests con `tusk test`. Cabe destacar que usamos el argumento `-coverprofile` para indicar a `go test` que realice un análisis de la cobertura y lo escriba en `coverage.txt`.

Por último, la fase `after_success` es la que se ejecuta si la *build* ha sido correcta. Lo que queremos hacer aquí es enviar el informe de cobertura generado anteriormente a Codecov, para lo que usamos su herramienta Bash.
