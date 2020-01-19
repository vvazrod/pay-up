---
layout: post
title: "Arquitectura del sistema"
---

Con PayUp, se pueden gestionar las deudas en grupos, de manera que sus integrantes añadan gastos y pagos al grupo y el sistema realice el cálculo de la deuda de forma automática.

A continuación, se van a presentar las entidades que componen el sistema, sus funcionalidades y el diseño de la arquitectura usada.

## Grupos (`gmicro`)

Un grupo está compuesto por una serie de integrantes, los cuales poseen un balance dentro del mismo. Este balance se altera cada vez que se añade un gasto o un pago en el grupo.

El microservicio que se va a implementar para esta entidad se encargará de la siguiente funcionalidad:

* Crear un grupo nuevo.
* Eliminar un grupo existente.
* Añadir integrantes a un grupo.
* Eliminar integrantes de un grupo.
* Calcular la deuda en el grupo.

## Gastos/Pagos (`tmicro`)

A un grupo se pueden añadir gastos realizados por alguno de sus integrantes para los demás miembros del grupo. Esto hace que se recalcule el balance de todos los integrantes. Para saldar deudas, un miembro del grupo debe realizar pagos a uno o varios de los demás.

Aunque los gastos y los pagos son de naturaleza diferente (los gastos son generales, mientras que los pagos son dirigidos), su funcionamiento es lo suficientemente similar como para poder considerarlos una única entidad, cuyo microservicio proporcionará la siguiente funcionalidad:

* Añadir un gasto.
* Añadir un pago.
* Eliminar el último gasto.
* Eliminar el último pago.

## Arquitectura

Las funcionalidades que hemos mencionado requieren de una respuesta por parte del cliente que las solicita, por lo que descartamos para nuestro sistema una arquitectura de paso de mensajes. Vamos a optar entonces por una arquitectura basada en microservicios, con uno para cada entidad descrita anteriormente.

Para tener un único punto de acceso a nuestro sistema por parte de los clientes, se usará un API *Gateway* basado en REST, el cual redirigirá las solicitudes a los distintos microservicios. La comunicación entre los distintos componentes se realizará con los protocolos HTTP/TCP y AMQP, usando JSON como formato de transmisión de datos.

![Diagrama de arquitectura del sistema](/pay-up/assets/images/architecture-diagram.png)

Como se puede ver en la imagen, usamos una cola de mensajes [RabbitMQ](https://www.rabbitmq.com/) para comunicar las peticiones al microservicio de gastos/pagos (`tmicro`). Esto se debe a que el orden en el que se procesen estas peticiones es importante para asegurar que el cálculo de los balances es el correcto. Una vez añadida la transacción a la base de datos, `tmicro` coloca un nuevo mensaje en la cola que debe procesar `gmicro` con el objetivo de actualizar el balance de las personas que intervienen en la transacción.

La implementación de los microservicios y del API *Gateway* se va a realizar con el lenguaje Go y haciendo uso del [Gorilla web toolkit](https://www.gorillatoolkit.org/).

Necesitaremos también un sistema de configuración remota, para lo que usaremos [Zookeeper](https://zookeeper.apache.org/). Para la persistencia de datos, utilizaremos bases de datos [PostgreSQL](https://www.postgresql.org/). Los logs producidos por cada parte del sistema deben ser centralizados, tarea que nos facilitará el paquete [Logrus](https://github.com/sirupsen/logrus).
