# PayUp

Con PayUp, los usuarios podrán crear grupos en los que añadir pagos realizados 
por uno o varios de los integrantes. El sistema automáticamente calculará la 
deuda y lo que debe pagar cada uno, recalculando esta cantidad si se añaden 
nuevos pagos. Así, se consigue simplificar los pagos y llevar un registro de 
los gastos.

Toda está funcionalidad será implementada como una serie de microservicios, los 
cuales se encargarán de:

* Autenticación de usuarios.
* Creación de grupos y adición de usuarios.
* Añadir gastos y pagos a un grupo.
* Redistribuir la deuda en un grupo.
* Notificar periódicamente a los usuarios sobre la deuda pendiente.

Para más información y documentación adicional de los hitos, visite la 
[página web del proyecto](https://varrrro.github.io/pay-up).
