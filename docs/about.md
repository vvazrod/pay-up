---
layout: page
title: Acerca del proyecto
permalink: /about/
---

Con PayUp, se pueden gestionar las deudas en grupos, de manera que sus integrantes
añadan gastos y pagos al grupo y el sistema realice el cálculo de la deuda de forma
automática.

A continuación, se van a presentar las entidades que componen el sistema, sus
funcionalidades y el diseño de la arquitectura usada.

# Grupos

Un grupo está compuesto por una serie de integrantes, los cuales poseen un balance
dentro del mismo. Este balance se altera cada vez que se añade un gasto o un pago en
el grupo.

El microservicio que se va a implementar para esta entidad se encargará de la siguiente
funcionalidad:

* Crear un grupo nuevo.
* Eliminar un grupo existente.
* Añadir integrantes a un grupo.
* Eliminar integrantes de un grupo.
* Calcular la deuda en el grupo.

# Gastos/Pagos

A un grupo se pueden añadir gastos realizados por alguno de sus integrantes para los
demás miembros del grupo. Esto hace que se recalcule el balance de todos los integrantes.
Para saldar deudas, un miembro del grupo debe realizar pagos a uno o varios de los demás.

Aunque los gastos y los pagos son de naturaleza diferente (los gastos son generales,
mientras que los pagos son dirigidos), su funcionamiento es lo suficientemente similar
como para poder considerarlos una única entidad, cuyo microservicio proporcionará la
siguiente funcionalidad:

* Añadir un gasto.
* Añadir un pago.
* Eliminar el último gasto.
* Eliminar el último pago.
