.. vim:set ft=rst spell:

Requirements for new Cloud Provider
===================================

This proposal should be lining out the necessary modifications to Taramk, to
make cloud-provider support more straight-forward.

Background
----------

While Tarmak was designed with the use-case of multiple providers in mind, it
currently only supports AWS  Part of this proposals is to investigate the
requirements for future cloud-provider implementations and how these could be
simplified.


Abstract requirements for a deployment of a Tarmak cluster
**********************************************************

Secure storage of arbitrary key/value pairs
+++++++++++++++++++++++++++++++++++++++++++

Used to auto unseal, other secrets are stored in terraform state

Metadata service
++++++++++++++++

.. todo::

    Not too sure where we use it, probably not neccessary as long our instances have
    an otherwise unique identifier. It came form you @james

VM Image repository
+++++++++++++++++++

(for storing base VM images etc)

 Packer uses it to generate base image

Docker Image repository
+++++++++++++++++++++++


Object storage
++++++++++++++

  TF state, backups


Load balancing
++++++++++++++

  Kubernetes master listens behind LB at AWS, scaling of master wouldn't allow for having a stable Hostname/IP otherwise

DNS
+++

  etcd nodes IP changes
  
  easier access of components

  vault "failover"
  

Private Network Peering
+++++++++++++++++++++++

  Multi cluster environments, need connectivity between hub <-- clusters (for accessing vault)

- Firewalling/L4 Access Control

  no node local firewall currently setup, relying on cloud provider firewalling, quite hard to do in software with changing IP addresses

.. todo::

   Provide some context to the situation now (e.g. current situation, short
   comings, future challenges, ...)


Objective
---------


Override sub methods
********************

- override provider areas in the environment (state buckets, dns, backup buckets, secret storage)
- requirements per cloud provider



.. todo::

    What are we going to do to make the situation described above better? (high
    level design overview, goals & considerations)


Changes
-------


Ensure configuration is done using Puppet for all roles
*******************************************************

vault/consul is done using bash


Notable items
-------------

Out of scope
------------

- actual implementations
