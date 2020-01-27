# PayUp

[![Build Status](https://travis-ci.com/Varrrro/pay-up.svg?branch=master)](https://travis-ci.com/Varrrro/pay-up)
[![codecov](https://codecov.io/gh/Varrrro/pay-up/branch/master/graph/badge.svg)](https://codecov.io/gh/Varrrro/pay-up)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

> **Nota:** Puede encontrar toda la documentación en la [página web del proyecto](https://varrrro.github.io/pay-up/).

## Herramienta de construcción

Para gestionar las distintas tareas relativas al proyecto, usaremos [Tusk](https://github.com/rliebz/tusk) como herramienta de construcción. Las tareas que se han definido se encuentran en el siguiente archivo:

> buildtool: tusk.yml

## Imágenes Docker

Para facilitar el despliegue de los distintos servicios en cualquier plataforma, se han contenerizado las distintas partes de la aplicación. Las imágenes creadas se pueden encontrar en el repositorio de Docker Hub que se indica a continuación.

> Contenedor: https://hub.docker.com/r/varrrro/pay-up

## Despliegue local

Puede desplegar el sistema en su ordenador usando Vagrant y Ansible con los siguientes comandos (__Nota:__ Tanto Vagrant como Ansible deben estar instalados en su ordenador):

```bash
cd deployments/vagrant
vagrant up
```

* Más información sobre el despliegue local con Vagrant [aquí](https://varrrro.github.io/pay-up/2020/01/27/local-deployment.html).
* Más información sobre el provisionamiento de las máquinas virtuales con Ansible [aquí](https://varrrro.github.io/pay-up/2020/01/27/provisioning.html).