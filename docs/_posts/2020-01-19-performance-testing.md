---
layout: post
title: "Evaluación de prestaciones"
---

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

## Resultados de las pruebas

Una vez definido el fichero de configuración de las pruebas, se procede a su ejecución. Se han realizado estas pruebas para distintos despliegues y, en todos los casos, se ha probado tanto con el `gateway` como sin él para comprobar su impacto en el rendimiento. A continuación, se van a presentar los resultados para cada despliegue.

### Despliegue con contenedores

Todas las partes del sistema a probar (`gateway`, `gmicro` y base de datos), se despliegan como contenedores usando las imágenes que hemos construido y la imagen oficial de PostgreSQL. Los resultados son los siguientes:

|Gateway|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Sí     |479,28 hits/s  |3,41% |18 ms          |34 ms         |121,16 KiB/s  |
|No     |527,35 hits/s  |10,75%|17 ms          |34 ms         |219,74 KiB/s  |

A parte de que no se alcanzan las prestaciones deseadas en ninguno de los dos casos, vemos como un porcentaje de las peticiones han acabado en error. Estos errores han sido todos del tipo `SocketException` y creemos que se deben a los recursos limitados de los que disponen los contenedores.

### Despliegue local con base de datos contenerizada

Ahora, vamos a ejecutar `gateway` y `gmicro` directamente en local sin contenedores, aunque sí que mantenemos la base de datos PostgreSQL en un contenedor. Los resultados son los siguientes:

|Gateway|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Sí     |756,62 hits/s  |0%    |11 ms          |24 ms         |133,73 KiB/s  |
|No     |828,28 hits/s  |0%    |10 ms          |23 ms         |146,39 KiB/s  |

Vemos como hemos eliminado los errores y las peticiones por segundo han mejorado considerablemente.

### Despliegue local completo

Dejamos los contenedores completamente de lado y montamos la base de datos PostgreSQL también en local junto con el `gateway` y `gmicro`, obteniendo los siguientes resultados:

|Gateway|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Sí     |945,05 hits/s  |0%    |9 ms           |18 ms         |167,03 KiB/s  |
|No     |1074,83 hits/s |0%    |8 ms           |17 ms         |189,97 KiB/s  |

Con este despliegue y realizando las pruebas directamente sobre `gmicro`, sin pasar por el `gateway`, sí que alcanzamos las prestaciones deseadas.

### Despliegue local con datos en memoria

Por último, vamos a hacer que `gmicro` devuelva directamente datos almacenados en memoria, sin acceder a la base de datos en ningún momento, para comprobar cuál es el impacto de esta en las prestaciones. Los resultados son los siguientes:

|Gateway|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Sí     |4090,42 hits/s |0%    |1 ms           |4 ms          |715,01 KiB/s  |
|No     |7249,57 hits/s |0%    |0 ms           |2 ms          |1269,76 KiB/s |

Como se puede ver, las prestaciones aumentan considerablemente en esta situación, lo que nos hace pensar que el cuello de botella de nuestro sistema se encuentra en el acceso a la base de datos.