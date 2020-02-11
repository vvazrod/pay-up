# PayUp

[![Build Status](https://travis-ci.com/Varrrro/pay-up.svg?branch=master)](https://travis-ci.com/Varrrro/pay-up)
[![codecov](https://codecov.io/gh/Varrrro/pay-up/branch/master/graph/badge.svg)](https://codecov.io/gh/Varrrro/pay-up)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

> __Note:__ This project was made for the Cloud Computing course of the Master's Degree in Computer Engineering at University of Granada. Developed with educational purposes only.

PayUp is a distributed system for managing payments and debts in a group. Create groups of people and start adding expenses, the system will calculate the balance of each member of the group so you don't have to.

You can find more info about the project and the technologies used in the [documentation](https://varrrro.github.io/pay-up/) (which is entirely in spanish, sorry friends :sweat_smile:)

## Running the system

You can deploy the complete system in three different ways:

* Locally, with Docker Compose: You need to have both Docker and Docker Compose installed on your computer. Then, run `docker-compose up` at `deployments/docker`.
* Locally, with Vagrant: You need to have both Vagrant and VirtualBox installed on your computer. Then, run `vagrant up` at `deployments/vagrant`. Ansible is also needed for provisioning.
* Remotely, with Ansible on GCP: Obviously, you need to have Ansible installed on your computer and also provide the credentials to your GCP project's service account. You can then run `ansible-playbook deploy.yml` at `deployments/ansible/gcp`.

In any of the above scenarios, you need to specify the environment variables that are needed by the deployment scripts.
