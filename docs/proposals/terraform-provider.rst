.. vim:set ft=rst spell:

Terraform Provider
==================

This proposal suggests how to approach the implementation of a terraform
provider for *Tarmak*, to make Tarmak <-> Terraform interactions within a
Terraform run more straightforward.

Background
----------

Right now the terraform code for the AWS provider (the only one implemented)
consists of multiple separate stacks (state, network, tools, vault,
kubernetes). Main reason for having these stacks is to enable Tarmak to do
operations in between parts of the resource spin up. Examples for such
operations are (*list might not be complete*):

- Bastion node needs to be up for other instances to check into wing
- Vault needs to be up and initialised, before PKI resources are created
- Vault needs to contain a clusters PKI resources, before kubernetes instances
  can be created (``init-tokens``)

The separation of stacks comes with some overhead for preparing terraform apply
(pull state, lock stack, plan run, apply run). Terraform can't make use of
parallel creation of resources that are independent from each other.


Objective
---------

An integration of these stacks into a single stack could lead to a substantial

reduction of execution time.

As terraform is running in a container is quite isolated from i

* Requires some terraform refactoring
* Should be done before implementing multiple providers

Changes
-------

Terraform code base
*******************

Terraform resources
*******************

A tarmak provider needs at least these 3 resources

``tarmak_bastion_instance``
~~~~~~~~~~~~~~~~~~~~~~~~~~~

A bastion instance

::

  Input:
  - bastion IP address or hostname
  - username for SSH

  Blocks till wing API server is healthy

``tarmak_vault_cluster``
~~~~~~~~~~~~~~~~~~~~~~~~

A vault cluster

::

  Input:
  - list of vault internal FQDNs or IP addresses

  Blocks till vault is initialised & unsealed


``tarmak_vault_instance_role``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This creates (once per role) an init token for such instances in Vault. 

::

  Input:
  - name of vault cluster
  - role name

  Output:
  - init_token per role


  Blocks till init token is setup



Notable items
-------------



Communication with the process in the container
***********************************************

I think one of the main difficulties is communication with the Tarmak process,
as terraform is run within a docker container with no communication available
to the main tarmak process (stdin/out is used by the terraform main process).

The proposal suggest that all ``terraform-provider-tarmak`` resources blocks
till the point, when the main Tarmak process connects using another exec to a
so called ``tarmak-connector`` executable that listens to a local Unix socket.
The ``terraform-provider-tarmak`` is then blocking until it is able to
establish a connection to the ``tarmak-connector`` on that socket.

This provides a secure and platform independent channel between Tarmak and
``terraform-provider-tarmak``.

::

   <<Tarmak>>  -- launches -- <<terraform|container>> 

       stdIN/OUT -- exec into  ---- <<exec terraform apply>>
                                    <<subprocess terraform-provider-tarmak
                                        |
                                     connects
                                        |
                                    unix socket
                                        |
                                     listens
                                        |
       stdIN/OUT -- exec into  ---- <<exec tarmak-connector>>


The protocol on that channel, should be using Golang's `net/rpc
<https://golang.org/pkg/net/rpc/>`_

PoC
***

I (`@simonswine <https://github.com/simonswine>`_) was starting a `PoC
<https://gitlab.jetstack.net/christian.simon/terraform-provider-tarmak/tree/master>`_
just to test how the terraform plugin model looks like.  It's not really having
anything implemented at this point, but might serve as a starting point.

Out of scope
------------

This proposal is not suggesting that we migrating features that are currently
done by the tarmak main process. The reason for that is that we don't want
terraform to become involved in the configuration provisioning of e..g the
vault cluster. This proposal should only improve the control we have from
Tarmak over things that happen in terraform.
