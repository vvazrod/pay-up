---
layout: post
title: "Configuración de Git"
---

Para este proyecto usamos Git y GitHub, por lo que lo primero que tuvimos
que hacer fue configurar correctamente estas herramientas. Para empezar,
se generó un par de claves SSH, las cuales se pueden ver en la imagen ya
generadas.

![Claves SSH generadas](/pay-up/assets/images/claves-ssh.png)

Subimos la clave pública `id_rsa.pub` a nuestra cuenta de GitHub para
permitir el acceso por SSH.

Por otra parte, también tenemos que configurar en Git nuestro nombre e
email para que aparezcan correctamente al realizar *commits*. Esta
configuración se realiza con los comandos de la siguiente imagen.

![Configuración de Git](/pay-up/assets/images/git-config.png)

Una vez hecho esto, ya estamos listos para crear nuestro repositorio
y empezar a trabajar en el proyecto.
