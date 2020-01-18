---
layout: post
title: "Gestión de procesos"
---

Cuando no lanzamos los servicios en un contenedor, necesitamos usar un gestor de procesos que permita controlar su ejecución fácilmente. Para ello, usamos una herramienta llamada [`supervisord`](http://supervisord.org/), con la cuál podemos definir la ejecución de cada servicio en un archivo de configuración.

A continuación, vamos a ver la composición de este fichero para el servicio `gmicro`.

```apacheconf
[supervisord]
nodaemon=true

[program:gmicro]
command=/usr/local/bin/gmicro

autostart=true
autorestart=true
startretries=10

stdout_logfile=program.log
stdout_logfile_maxbytes=0
```

Bajo la etiqueta `[supervisord]` se definen las características del propio gestor de procesos para esta configuración. Lo único que indicamos aquí es que el gestor se ejecute en primer plano con `nodaemon`.

En `[program:gmicro]` es dónde definimos la configuración para la aplicación concreta. Con `command` se establece el comando que debe ejecutar el gestor, en nuestro caso, solo debe ejecutar el binario. Los parámetros `autostart` y `autorestart` los fijamos a _true_ para indicar que la aplicación se debe lanzar automáticamente al iniciar el gestor de procesos y que se debe relanzar si su ejecución se termina por alguna causa. Con `startretries` se establece el número de intentos. Finalmente, indicamos la ruta del fichero en el que se escribiran los logs de `stdout` de la aplicación con `stdout_logfile`.

Para ejecutar el servicio `gmicro` usando esta configuración, primero debemos asegurarnos que el ejecutable se encuentra en `/usr/local/bin` y luego ejecutar `supervisord -c` indicando la ruta del archivo de configuración.
