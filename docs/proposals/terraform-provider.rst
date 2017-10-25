Terraform Provider
------------------

Motivations
***********

Right now the terraform code for the AWS provider consists for multiple
separate stacks (state, network, tools, vault, kubernetes). Main reason for the
the different stacks was to enable operations in between parts of the resource
spin up. These include (list might not be complete):

- Bastion node needs to be up for other instances to check into wing
- Vault needs to be up and initialised, before PKI resources are created
- Vault needs to contain a clusters PKI resources, before kubernetes instances
  can be created (``init-tokens``)

Design ideas
************

A tarmak provider needs at least these 3 resources

``tarmak_bastion_instance``
~~~~~~~~~~~~~~~~~~~~~~~~~~~

A bastion instance

::

  Input:
  - bastion IP address or hostname
  - username for SSH

  Blocks till wing API server is running

``tarmak_vault_cluster``
~~~~~~~~~~~~~~~~~~~~~~~~

A vault cluster

::

  Input:
  - list of vault internal FQDN or IP addresses

  Blocks till vault is ready


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


Difficulties
************

I think the main difficulty is communication with the Tarmak process, as
terraform is run within a docker container with no communication available to
the main tarmak process (stdin/out is used by the terraform main process).

The proposal suggest that all ``terraform-provider-tarmak`` resources blocks
till the point, when the main Tarmak process connects using another exec to a
so called ``tarmak-connector`` executable that speaks via a local Unix socket
to the ``terraform-provider-tarmak``.

This provides a secure and platform independent channel between Tarmak and
``terraform-provider-tarmak``.

::

   <<Tarmak>>  -- launches -- <<terraform|container>> 

       stdIN/OUT -- exec into  ---- <<exec terraform apply>>
                                    <<subprocess terraform-provider-tarmak
                                        |
                                    unix socket
                                        |
       stdIN/OUT -- exec into  ---- <<exec tarmak-connector>>


The protocol on that channel, should be either an Unix socket communication
using Kubernetes API features or `net/rpc <https://golang.org/pkg/net/rpc/>`_

PoC
***

I (`@simonswine <https://github.com/simonswine>`_) was starting a `PoC
<https://gitlab.jetstack.net/christian.simon/terraform-provider-tarmak/tree/master>`_
just to test how the terraform plugin model looks like.  It's not really having
anything implemented at this point, but might serve as a starting point.
