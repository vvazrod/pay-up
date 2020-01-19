---
layout: post
title: "Evaluación de prestaciones"
---

Una vez definido el fichero de configuración de las pruebas, se procede a su ejecución. Se han realizado estas pruebas para distintos despliegues y, en todos los casos, se ha probado tanto con el `gateway` como sin él para comprobar su impacto en el rendimiento. A continuación, se van a presentar los resultados para cada despliegue.

## Despliegue con contenedores

Todas las partes del sistema a probar (`gateway`, `gmicro` y base de datos), se despliegan como contenedores usando las imágenes que hemos construido y la imagen oficial de PostgreSQL. Los resultados son los siguientes:

|Gateway|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Sí     |479,28 hits/s  |3,41% |18 ms          |34 ms         |121,16 KiB/s  |
|No     |527,35 hits/s  |10,75%|17 ms          |34 ms         |219,74 KiB/s  |

A parte de que no se alcanzan las prestaciones deseadas en ninguno de los dos casos, vemos como un porcentaje de las peticiones han acabado en error. Estos errores han sido todos del tipo `SocketException` y creemos que se deben a los recursos limitados de los que disponen los contenedores.

## Despliegue local con base de datos contenerizada

Ahora, vamos a ejecutar `gateway` y `gmicro` directamente en local sin contenedores, aunque sí que mantenemos la base de datos PostgreSQL en un contenedor. Los resultados son los siguientes:

|Gateway|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Sí     |756,62 hits/s  |0%    |11 ms          |24 ms         |133,73 KiB/s  |
|No     |828,28 hits/s  |0%    |10 ms          |23 ms         |146,39 KiB/s  |

Vemos como hemos eliminado los errores y las peticiones por segundo han mejorado considerablemente.

## Despliegue local completo

Dejamos los contenedores completamente de lado y montamos la base de datos PostgreSQL también en local junto con el `gateway` y `gmicro`, obteniendo los siguientes resultados:

|Gateway|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Sí     |945,05 hits/s  |0%    |9 ms           |18 ms         |167,03 KiB/s  |
|No     |1074,83 hits/s |0%    |8 ms           |17 ms         |189,97 KiB/s  |

Con este despliegue y realizando las pruebas directamente sobre `gmicro`, sin pasar por el `gateway`, sí que alcanzamos las prestaciones deseadas.

## Despliegue local con datos en memoria

Por último, vamos a hacer que `gmicro` devuelva directamente datos almacenados en memoria, sin acceder a la base de datos en ningún momento, para comprobar cuál es el impacto de esta en las prestaciones. Los resultados son los siguientes:

|Gateway|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Sí     |4090,42 hits/s |0%    |1 ms           |4 ms          |715,01 KiB/s  |
|No     |7249,57 hits/s |0%    |0 ms           |2 ms          |1269,76 KiB/s |

Como se puede ver, las prestaciones aumentan considerablemente en esta situación, lo que nos hace pensar que el cuello de botella de nuestro sistema se encuentra en el acceso a la base de datos.