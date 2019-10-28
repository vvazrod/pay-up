# PayUp

La arquitectura del sistema definida en el hito anterior se puede encontrar en [este *post*](https://varrrro.github.io/pay-up/2019/10/28/system-architecture.html). Aunque ahí ya indicamos algunas de las herramientas que íbamos a utilizar, a continuación vamos a describirlas en mayor profundidad.

Para la implementación, usaremos el *toolkit* [Go-kit](https://github.com/go-kit/kit), ayudándonos del [paquete mux de Gorilla web toolkit](https://github.com/gorilla/mux) para el enrutamiento de peticiones en el API Gateway.

La caché de nuestro sistema se realizará con [Memcached](https://memcached.org/), usando el paquete [gomemcache](https://github.com/bradfitz/gomemcache) para trabajar con esta desde nuestro código.

El almacenamiento persistente correrá a cuenta de [PostgreSQL](https://www.postgresql.org/), usando como ORM el paquete [pg](https://github.com/go-pg/pg).

Usaremos [Zookeeper](https://zookeeper.apache.org/) para la configuración remota, con el cliente nativo de Go [go-zookeeper](https://github.com/samuel/go-zookeeper).

La conexión con el sistema de logs [Graylog](https://www.graylog.org/) se realizará con su formato propio GELF, haciendo uso del paquete [go-gelf](https://github.com/Graylog2/go-gelf).

Para comunicarnos con la cola de mensajes [RabbitMQ](https://www.rabbitmq.com/), usaremos el protocolo AMQP con el paquete [amqp](https://github.com/streadway/amqp).
