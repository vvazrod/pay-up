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

## Evaluación de prestaciones

Se han realizado pruebas de carga sobre el sistema usando [Taurus](https://gettaurus.org/) para evaluar su nivel de prestaciones. El archivo de evaluación es el siguiente:

> Prestaciones: load_test.yml

Puede leer más sobre estas pruebas [aquí](https://varrrro.github.io/pay-up/2020/01/19/performance-testing.html).

## Despliegue

Puede desplegar el sistema en su ordenador usando Vagrant y Ansible con los siguientes comandos (__Nota:__ Vagrant, VirtualBox y Ansible deben estar instalados en su ordenador):

```bash
cd deployments/vagrant
vagrant up
```

Por otra parte, también puede desplegar el sistema en _Google Cloud Platform_ añadiendo un fichero `credentials.json` al directorio `deployments/ansible/gcp` con los datos de acceso de una cuenta de servicio a un proyecto en GCP. También debe incluir el usuario de su clave SSH (`ssh_user`) y su contraseña (`ssh_pass`) en un fichero `ssh_credentials.json` ubicado en `deployments/ansible/env/`. Una vez hecho esto, use los siguientes comandos:

```bash
cd deployments/ansible/gcp
ansible-playbook deploy.yml
```

Ó, si está usando Python 3:

```bash
cd deployments/ansible/gcp
ansible-playbook deploy.yml -e'ansible_python_interpreter=/usr/bin/python3'
```

Tanto si va a realizar el despliegue en local como en remoto, debe añadir también a la carpeta `provision/env` los ficheros `gateway.yml`, `gmicro.yml` y `tmicro.yml` con las variables de entorno necesarias para cada servicio.

Puede obtener más información sobre el despliegue y provisionamiento de máquinas virtuales [aquí](https://varrrro.github.io/pay-up/2020/01/27/deployment-and-provisioning.html).