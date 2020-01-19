# PayUp

[![Build Status](https://travis-ci.com/Varrrro/pay-up.svg?branch=master)](https://travis-ci.com/Varrrro/pay-up)
[![codecov](https://codecov.io/gh/Varrrro/pay-up/branch/master/graph/badge.svg)](https://codecov.io/gh/Varrrro/pay-up)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

> **Nota:** Puede encontrar toda la documentación en la [página web del proyecto](https://varrrro.github.io/pay-up/).

## Estado del proyecto

Actualmente, se ha implementado toda la funcionalidad pensada para el sistema, el cual se compone de `gmicro`, `tmicro` y el `gateway` (puede leer más sobre la arquitectura del sistema [aquí](https://varrrro.github.io/pay-up/2019/10/28/system-architecture.html). A continuación, se van a presentar los aspectos más importantes de las adiciones.

### Inversión de dependencias

Para no unir nuestra implementación a un gestor de base de datos concreto, se hace uso de un ORM llamado `gorm`. Usando este ORM, definimos el modelo de datos usando `struct` de Go y trabajamos con una conexión de base de datos genérica proporcionada por el paquete, que nos permite realizar todas las operaciones CRUD (_Create, Read, Update, Delete_).

Tanto en `gmicro` como en `tmicro`, implementamos un _data manager_ que hace uso de esta conexión y que sirve como única fuente de verdad, es decir, se encarga de realizar todas las operaciones necesarias de acceso al modelo de datos.

En la función `main` del servicio, se crea la conexión usando una base de datos concreta y con la que luego se construye el _data manager_, pero nuestra implementación es ajena a esta elección.

### Cola de mensajes

Se usa una cola de mensajes RabbitMQ para las peticiones relacionadas con los gastos y los pagos, asegurando que se procesen en el orden de llegada. Tanto `tmicro` como `gmicro` trabajan con un mismo _exchange_ pero con una cola o _queue_ para cada uno, de forma que ambos lanzan un consumidor AMQP con una _goroutine_ que se encarga de enviar los mensajes recibidos en la cola correspondiente a un controlador.

## Evaluación de prestaciones

El sistema desarrollado debe estar preparado para soportar una determinada carga, que marcamos en las 1000 peticiones por segundo contando con 10 usuarios concurrentes y sin que se produzcan errores. Vamos a comprobar que se cumplan estas prestaciones mínimas usando la herramienta [Taurus](https://gettaurus.org/).

Para la comprobación de las prestaciones, se va a trabajar solo con `gmicro` y el `gateway`, ya que las peticiones destinadas a `tmicro` no dan una respuesta al usuario después de ser procesadas, si no que es el `gateway` el que devuelve una respuesta de operación aceptada al comprobar que su formato es correcto.

Taurus nos permite definir la configuración de las pruebas de prestaciones en un fichero YAML.

> Prestaciones: load_test.yml

Vamos a presentar este fichero y explicar por partes su composición.

```yaml
execution:
  - scenario: gmicro
    concurrency: 10
    ramp-up: 10s
    hold-for: 50s
```
    
La etiqueta `execution` se usa para presentar la configuración de la ejecución de las pruebas. Aquí se pueden fijar propiedades globales y se definen los escenarios de prueba que se van a ejecutar. En nuestro caso, definimos un único escenario llamado _gmicro_, el cuál va a tener un total de 10 hilos o _threads_ concurrentes (etiqueta `concurrency`), los cuales representan nuestros 10 usuarios ficticios. El número de hilos irá aumentando progresivamente hasta llegar a los 10 en 10 segundos (`ramp-up`) y los mantendrá todos juntos durante otros 50 segundos (`hold-for`), con lo que las pruebas durarán 1 minuto en total.

```yaml
scenarios:
  gmicro:
    variables:
      id: ${__UUID}

    requests:
      - once:
        - label: Create group
          url: http://localhost:8080/groups
          method: POST
          headers:
            Content-Type: application/json
          body: '{"id":"${id}", "name":"test"}'

      - label: Fetch group
        url: http://localhost:8080/groups/${id}
        method: GET

```

Como su propio nombre indica, la etiqueta `scenarios` la usamos para definir las pruebas que conforman los distintos escenarios (nosotros solo usamos uno). Las pruebas van a ser sencillas, primero se va a crear un grupo en el microservicio y luego se van a realizar peticiones para pedir los datos de este grupo nuevo. Estas peticiones se definen en el apartado `requests`.

El orden de ejecución será el mismo que el de definición de las peticiones. Para cada una, proporcionamos un `url` y un verbo HTTP (`method`). En el caso de la petición de creación del grupo, debemos fijar el _header_ `Content-Type` a `application/json` dentro de la sección `headers`. También definimos el `body` de esta petición POST. Con la etiqueta `label` podemos definir un nombre para cada petición.

Para cada hilo, queremos que se ejecute solo una vez la petición POST y luego se realicen multitud de peticiones GET sobre el recurso creado. Con la etiqueta `once` hacemos que la petición de creación del recurso se realice solo una vez por hilo.

En nuestro sistema, es el cliente el que genera los UUIDs de los recursos que se van a crear (el servidor solo realiza comprobaciones de unicidad y formato). Es por ello que definimos la variable `id` en el escenario. Usando la función de JMeter `${__UUID}`, generamos un nuevo UUID en cada hilo y lo guardamos en esta variable. Como se puede ver, se utiliza esta variable en las peticiones para que cada hilo realice peticiones sobre un recurso independiente creado por si mismo.

Una vez definida la configuración de las pruebas, se procede a su realización con distintos despliegues y al análisis de los resultados. Puede leer más sobre estos resultados [aquí](https://varrrro.github.io/pay-up/2020/01/19/performance-testing.html).

## Herramienta de construcción

Para gestionar las distintas tareas relativas al proyecto, usaremos [Tusk](https://github.com/rliebz/tusk) como herramienta de construcción. Las tareas que se han definido se encuentran en el siguiente archivo:

> buildtool: tusk.yml

## Imágenes Docker

Para facilitar el despliegue de los distintos servicios en cualquier plataforma, se han contenerizado las distintas partes de la aplicación. Las imágenes creadas se pueden encontrar en el repositorio de Docker Hub que se indica a continuación.

> Contenedor: https://hub.docker.com/r/varrrro/pay-up
