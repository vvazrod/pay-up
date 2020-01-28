---
layout: post
title: "Despliegue y provisionamiento"
---

Para el despliegue de nuestro proyecto, hemos definido la infraestructura de máquinas virtuales necesaria como código tanto para el despliegue local como remoto usando Vagrant y Ansible, respectivamente. Estas máquinas virtuales, a su vez, se provisionan de todo lo necesario para la ejecución de los servicios usando Ansible.

En las siguientes secciones vamos a explicar en profundidad como se ha conseguido esto, además de evaluar el rendimiento del sistema en ambos despliegues.

## Infraestructura a desplegar

Para desplegar todos nuestros servicios, hemos decidido usar un total de 3 máquinas virtuales, las cuales se van a organizar así:

* VM #1: Aquí se ejecutarán tanto el servicio `gmicro` como su base de datos, además de la cola de mensajes RabbitMQ del sistema.
* VM #2: El servicio `tmicro` y su base de datos se ejecutarán en esta máquina.
* VM #3: La máquina más pequeña, dará cabida tan solo al `gateway`.

La razón para esta organización es que, de esta forma, las máquinas #2 y #3 hacen referencia a la #1 y se simplifican las conexiones de red. Además, se decide mantener en la misma máquina un servicio y su base de datos correspondiente para reducir la latencia y agilizar las peticiones.

Las tres máquinas virtuales deberán formar parte de la misma red privada, teniendo la máquina #1 una IP interna estática que usarán #2 y #3 para conectarse a ella (usaremos la misma IP tanto en el despliegue local como el remoto) y #3 una IP pública que permita el uso del sistema.

## Provisionamiento

El provisionamiento de las máquinas virtuales se va a realizar mediante unos _playbooks_ de Ansible que son utilizados indistintamente por ambos despliegues. Tendremos un _playbook_ específico para cada máquina que hemos definido en la sección anterior, además de un rol `common` que hemos creado para agrupar las tareas comunes a todas las máquinas. Vamos a comenzar viendo el fichero de tareas de este rol.

```yml
---
- name: Install docker
  include_role:
    name: geerlingguy.docker
  vars:
    docker_install_compose: false

- name: Install pip3
  apt:
    pkg: python3-pip
    update_cache: yes

- name: Install docker module
  pip:
    name: docker
```

El despliegue del software dentro de las máquinas (servicios, bases de datos, cola AMQP) se realizará con contenedores Docker debido a su facilidad de uso. Por ello, las tareas comunes a todas las máquinas consisten básicamente en instalar Docker haciendo uso de otro rol `geerlingguy.docker` de _Ansible Galaxy_ e instalar `pip` para instalar a continuación el módulo de Python que usa Ansible para trabajar con Docker. Una vez visto el rol `common`, vamos a ver como se aplica a uno de los _playbooks_ específicos de las máquinas, en concreto, al de la #1.

```yml
---
- hosts: gmicro
  become: yes
  vars_files:
    - env/gmicro.yml
  roles:
    - common
```

Para no confundirnos, hemos apodado cada máquina con el nombre del servicio que acogen. Por eso, este _playbook_ actúa sobre los _hosts_ `gmicro`. Hacemos uso de `become` para ejecutar las tareas como superusuario, ya que nos será necesario. En el apartado `vars_files` indicamos la ruta al fichero que contiene las variables de entorno necesarias para esta máquina (valores de conexión a la cola AMQP, usuario y contraseña de acceso a la base de datos, ...) y, en el apartado `roles`, añadimos nuestro rol `common`.

Vamos a ver ahora las tareas.

```yml
tasks:
    - name: Create docker network
      docker_network:
        name: main

    - name: Run RabbitMQ container
      docker_container:
        name: rabbit
        image: rabbitmq:3
        detach: yes
        networks:
          - name: main
        purge_networks: yes
        ports:
          - "5672:5672"

    - name: Run PostgreSQL container
      docker_container:
        name: db-gmicro
        image: postgres:12
        detach: yes
        networks:
          - name: main
        purge_networks: yes
        env:
          POSTGRES_USER: "{{ db_user }}"
          POSTGRES_PASSWORD: "{{ db_pass }}"
          POSTGRES_DB: "{{ db_name }}"

    - name: Run gmicro container
      docker_container:
        name: gmicro
        image: varrrro/pay-up:gmicro
        detach: yes
        networks:
          - name: main
        purge_networks: yes
        ports:
          - "8080:8080"
        env:
          RABBIT_CONN: "{{ rabbit_conn }}"
          DB_TYPE: "{{ db_type }}"
          DB_CONN: "{{ db_conn }}"
          EXCHANGE: "{{ exchange }}"
          QUEUE: "{{ queue }}"
          CTAG: "{{ ctag }}"
```

Al tener varios contenedores que deben comunicarse entre sí, lo primero que hacemos es crear una red Docker con `docker_network`. Después, arrancamos cada uno de los contenedores con `docker_container` y algunos parámetros:

* `name` e `image` son, como se puede suponer, el nombre del nuevo contenedor y la imagen a utilizar, respectivamente.
* Fijando `detach` a _yes_ hacemos que el contenedor se ejecute en segundo plano, permitiendo que continúe la ejecución de las tareas.
* Con `networks`, conectamos el nuevo contenedor a la red que creamos al principio. Debemos activar también `purge_networks` para que elimine las redes por defecto y deje solo la nuestra, ya que si no da problemas.
* Usamos `ports` para enlazar puertos de la máquina virtual a los puertos del contenedor que nos interesan (en este caso, 8080 del servicio para mandarle peticiones y 5672 de la cola para conectar con ella).
* `env` nos permite definir una serie de variables de entorno con las que se va a ejecutar el contenedor. Aquí es donde usamos los valores del fichero `env/gmicro.yml` definido anteriormente.

Una vez terminadas estas tareas, la máquina virtual tiene todos los procesos necesarios ejecutándose y está lista para ser utilizada.

Los _playbooks_ para las máquinas de `tmicro` y `gateway` son análogos al que hemos visto.

## Despliegue local

Para el desarrollo, vamos a configurar un despliegue automático de nuestra infraestructura de máquinas virtuales en nuestra propia máquina física usando Vagrant. Esta herramienta nos va a permitir definir las máquinas virtuales y su provisionamiento en un fichero `Vagrantfile` y lanzarlas de manera sencilla. Vamos a echar un vistazo a este fichero.

```ruby
Vagrant.configure("2") do |config|

  # VM for the gmicro microservice, its database and the AMQP queue.
  config.vm.define "gmicro" do |gmicro|
    gmicro.vm.box = "ubuntu/bionic64"

    # Set private static IP address.
    gmicro.vm.network "private_network", ip: "10.0.0.10"

    # Use VirtualBox to create this VM with 2 cores and 4GB of RAM.
    gmicro.vm.provider "virtualbox" do |vb|
      vb.cpus = 2
      vb.memory = "4096"
    end

    # Provision VM with Ansible.
    gmicro.vm.provision "ansible" do |ansible|
      ansible.playbook = "../../provision/gmicro.yml"
    end
  end

  # VM for the tmicro microservice and its database.
  config.vm.define "tmicro" do |tmicro|
    tmicro.vm.box = "ubuntu/bionic64"

    # Set private IP address using DHCP.
    tmicro.vm.network "private_network", type: "dhcp"

    # Use VirtualBox to create this VM with 2 cores and 4GB of RAM.
    tmicro.vm.provider "virtualbox" do |vb|
      vb.cpus = 2
      vb.memory = "4096"
    end

    # Provision VM with Ansible.
    tmicro.vm.provision "ansible" do |ansible|
      ansible.playbook = "../../provision/tmicro.yml"
    end
  end

  # VM for the API gateway.
  config.vm.define "gateway" do |gateway|
    gateway.vm.box = "ubuntu/bionic64"

    # Set private IP address using DHCP.
    gateway.vm.network "private_network", type: "dhcp"

    # Forward port to access system from host.
    gateway.vm.network "forwarded_port", guest: 8080, host: 8080 

    # Use VirtualBox to create this VM with 1 core and 2GB of RAM.
    gateway.vm.provider "virtualbox" do |vb|
      vb.cpus = 1
      vb.memory = "2048"
    end

    # Provision VM with Ansible.
    gateway.vm.provision "ansible" do |ansible|
      ansible.playbook = "../../provision/gateway.yml"
    end
  end
end
```

De aquí, vamos a destacar algunas cosas:

* Definimos la IP interna estática para la máquina de `gmicro` con `gmicro.vm.network "private_network", ip: "10.0.0.10"`. Como la IP del resto de máquinas no necesita ser estática, usamos `type: "dhcp"`.
* Enlazamos el puerto 8080 de la máquina física con el 8080 de la máquina virtual de `gateway` con `gateway.vm.network "forwarded_port", guest: 8080, host: 8080 ` para tener así acceso al sistema.
* Con `gateway.vm.provision "ansible"` definimos que vamos a usar Ansible como provisionador e indicamos la ruta al _playbook_ a ejecutar para cada máquina.

Como detalle importante, al enlazar Ansible como provisionador directamente en el `Vagrantfile`, el propio Vagrant se encarga de gestionar la conexión SSH de Ansible con la máquina virtual, permitiéndonos olvidarnos del fichero de ivnentario típico de Ansible.

Si usamos `vagrant up`, lanzaremos las tres máquinas virtuales y desplegaremos en ellas el sistema, el cuál queda totalmente funcional y accesible a través de `localhost:8080`.

## Despliegue remoto

Para el despliegue remoto, hemos elegido _Google Cloud Platform_ como proveedor cloud y Ansible como herramienta para la definición de la infraestructura virtual, además de su provisionamiento.

Lo primero que debemos hacer es crear un proyecto en el _Google Cloud Console_ y añadir nuestra clave SSH pública al mismo. Esta clave se copiará automáticamente a todas las instancias de máquinas virtuales que se lancen, permitiendo a Ansible conectarse a ellas y realizar el provisionamiento.

Para poder crear la infraestructura desde Ansible necesitaremos, además, crear una cuenta de servicio con permisos de edición y descargar el fichero JSON de credenciales.

Vamos a explica ahora, por partes, el contenido del _playbook_ `deploy.yml`, que es el encargado de realizar este despliegue remoto.

```yml
---
- name: Create GCP infrastructure
  hosts: localhost
  gather_facts: no

  vars_files:
    - ../env/ssh_credentials.yml

  vars:
    service_account_file: ./credentials.json
    project: payup-2020
    auth_kind: serviceaccount
    region: "europe-west1"
    zone: "europe-west1-b"
    scopes:
      - https://www.googleapis.com/auth/compute
```

Para empezar, debemos usar `localhost` como _host_ objetivo ya que Ansible no permite ejecutar un _playbook_ sin definir uno. En `vars_files` añadimos el fichero que contiene el usuario SSH de la clave que hemos añadido al proyecto y su contraseña (`ssh_user` y `ssh_pass`), necesarios para el provisionamiento de las máquinas.

Luego, hay varias variables que debemos definir en `vars`, como son la ruta del fichero de credenciales de la cuenta de servicio, el ID del proyecto de GCP, la región o la zona. Estas variables se usarán a lo largo del _playbook_ para crear los recursos.

```yml
    - name: Create network
      gcp_compute_network:
        name: "payup-net"
        auto_create_subnetworks: no
        project: "{{ project }}"
        auth_kind: "{{ auth_kind }}"
        service_account_file: "{{ service_account_file }}"
        scopes: "{{ scopes }}"
      register: network

    - name: Create firewall
      gcp_compute_firewall:
        name: "payup-firewall"
        allowed:
        - ip_protocol: tcp
          ports:
          - "22"
          - "8080"
          - "5672"
        network: "{{ network }}"
        project: "{{ project }}"
        auth_kind: "{{ auth_kind }}"
        service_account_file: "{{ service_account_file }}"
        scopes: "{{ scopes }}"

    - name: Create subnet
      gcp_compute_subnetwork:
        name: "payup-subnet-eu"
        region: "{{ region }}"
        network: "{{ network }}"
        ip_cidr_range: 10.0.0.0/8
        project: "{{ project }}"
        auth_kind: "{{ auth_kind }}"
        service_account_file: "{{ service_account_file }}"
        scopes: "{{ scopes }}"
      register: subnet
```

Lo primero que hacemos es definir la infraestructura de red: creamos una nueva VPC con `gcp_compute_network`, le aplicamos un _firewall_ con `gcp_compute_firewall` y creamos una subred con `gcp_compute_subnetwork`. El _firewall_ es necesario para permitir el tráfico TCP relativo a los puertos 22 (SSH), 8080 (peticiones HTTP a las APIs) y 5672 (conexiones a la cola AMQP). Por otra parte, en la subred indicamos un rango pequeño de direcciones ya que solo tendremos 3 máquinas en ella.

```yml
    - name: Reserve static internal IP address for gmicro
      gcp_compute_address:
        name: "gmicro-internal-ip"
        address: 10.0.0.10
        address_type: "INTERNAL"
        subnetwork: "{{ subnet }}"
        region: "{{ region }}"
        project: "{{ project }}"
        auth_kind: "{{ auth_kind }}"
        service_account_file: "{{ service_account_file }}"
        scopes: "{{ scopes }}"
      register: gmicro_internal_ip

    - name: Reserve external IP address for gmicro
      gcp_compute_address:
        name: "gmicro-external-ip"
        address_type: "EXTERNAL"
        region: "{{ region }}"
        project: "{{ project }}"
        auth_kind: "{{ auth_kind }}"
        service_account_file: "{{ service_account_file }}"
        scopes: "{{ scopes }}"
      register: gmicro_external_ip
```

Ahora, reservamos tanto la IP estática interna como la IP externa para `gmicro`. Todas las máquinas remotas deben disponer de una IP pública para que Ansible se pueda conectar a ellas y realizar el provisionamiento, por lo que este último paso lo repetimos para cada máquina.

El principal detalle es que volvemos a asignar la misma dirección que en el despliegue local, `10.0.0.10`, lo que nos permite poder usar las mismas variables de conexiones entre servicios en nuestros _playbooks_.

Con `register`, vamos guardando los recursos generados en variables para poder utilizarlos en otras partes del _playbook_.

```yml
    - name: Create disk for gmicro
      gcp_compute_disk:
        name: "gmicro-disk"
        size_gb: 50
        source_image: projects/ubuntu-os-cloud/global/images/ubuntu-1804-bionic-v20200108
        zone: "{{ zone }}"
        project: "{{ project }}"
        auth_kind: "{{ auth_kind }}"
        service_account_file: "{{ service_account_file }}"
        scopes: "{{ scopes }}"
      register: gmicro_disk

    - name: Create gmicro VM
      gcp_compute_instance:
        name: "gmicro-vm"
        machine_type: n1-standard-1
        disks:
          - auto_delete: true
            boot: true
            source: "{{ gmicro_disk }}"
        network_interfaces:
          - network: "{{ network }}"
            subnetwork: "{{ subnet }}"
            network_ip: "{{ gmicro_internal_ip.address }}"
            access_configs:
            - name: External NAT
              nat_ip: "{{ gmicro_external_ip }}"
              type: ONE_TO_ONE_NAT
        zone: "{{ zone }}"
        project: "{{ project }}"
        auth_kind: "{{ auth_kind }}"
        service_account_file: "{{ service_account_file }}"
        scopes: "{{ scopes }}"
      register: gmicro_vm
```

Creamos el disco y, con ese disco, la máquina virtual (una vez más, paso que repetimos para cada una de las máquinas). En la creación del disco indicamos su tamaño (`size_gb`) y la imagen de sistema operativo que se va a usar (`source_image`), mientras que en la creación de la máquina virtual indicamos que este disco es de arranque (`boot`) y definimos todas las interfaces (`network_interfaces`).

```yml
    - name: Wait for SSH to come up on gmicro VM
      wait_for: host={{ gmicro_external_ip.address }} port=22 delay=10 timeout=60

    - name: Add gmicro host
      add_host:
        hostname: "{{ gmicro_external_ip.address }}"
        groupname: gmicro
        ansible_user: "{{ ssh_user }}"
        ansible_password: "{{ ssh_pass }}"
```

Debido a que las máquinas tardan en arrancar y estar listas, debemos esperar a que el puerto 22 esté activo antes de poder continuar con el provisionamiento de las mismas. Esto lo conseguimos con `wait_for`.

Una vez que el puerto está listo, añadimos la máquina al inventario dinámico de _hosts_ que Ansible mantiene en memoria durante la ejecución (`add_hosts`). Aquí, usamos las credenciales SSH del fichero definido anteriormente.

```yml
- name: Provision gmicro VM
  import_playbook: ../../../provision/gmicro.yml
```

Por último, ejecutamos el _playbook_ de provisionamiento de la máquina para cada una de ellas. Al final, obtenemos un sistema completamente funcional al que nos podemos conectar con la IP pública que genera GCP, que podemos mostrar por consola así:

```yml
    - name: Show gateway public IP
      debug:
        msg: "Gateway's public IP is {{ gateway_external_ip.address }}"
```

## Comparación de rendimiento

Para evaluar el rendimiento de nuestro sistema en cada uno de los despliegues, se han ejecutado las mismas pruebas realizadas para la evaluación de prestaciones, obteniendo los siguientes resultados:

|Despliegue|Avg. Throughput|Errors|Avg. Resp. time|90% Resp. time|Avg. Bandwidth|
|-------|---------------|------|---------------|--------------|--------------|
|Remoto     |116,08 hits/s  |0% |79 ms          |89 ms         |20,51 KiB/s  |
|Local     |419,13 hits/s  |0%|21 ms          |61 ms         |74,07 KiB/s  |

Como podemos ver, el rendimiento cae enormemente al realizar estos despliegues. Hay que tener en cuenta que se han realizado con contenedores dentro de las máquinas virtuales, siendo esta configuración la que peor rendimiento nos daba en las pruebas realizadas en local, por lo que podemos suponer que si realizamos el despliegue directamente sobre las máquinas, sin utilizar contenedores, mejoraríamos el rendimiento.

No obstante, se obtienen resultados mucho peores también a los obtenidos realizando este despliegue con contenedores en nuestra máquina física, debido esencialmente a la latencia añadida.