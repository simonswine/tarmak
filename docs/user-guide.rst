.. getting-started:

User guide
==========

Getting started with AWS
------------------------

In this getting started guide, we walk through how to use initialise Tarmak with a new Provider (AWS) and Environment and provision a Kubernetes cluster.
This will comprise Kubernetes master and worker nodes, etcd clusters, Vault and a bastion node with a public IP address
(see :ref:`Architecture overview <architecture_overview>` for details of cluster components)

Prerequisites
~~~~~~~~~~~~~

* Docker
* An AWS account
* A public DNS zone that can be delegated to AWS Route 53
* Optional: Vault with the `AWS secret backend <https://www.vaultproject.io/docs/secrets/aws/index.html>`_ configured

Overview of steps to follow
~~~~~~~~~~~~~~~~~~~~~~~~~~~

* :ref:`Initialise cluster configuration <init_config>`
* :ref:`Build an image (AMI) <create_ami>`
* :ref:`Create the cluster <create_cluster>`
* :ref:`Use the cluster <use_cluster>`
* :ref:`Destroy the cluster <destroy_cluster>`

.. _init_config:

Initialise configuration
~~~~~~~~~~~~~~~~~~~~~~~~

Simply run ``tarmak init`` to initialise configuration for the first time. You will be prompted for the necessary configuration
to set-up a new :ref:`Provider <providers_resource>` (AWS) and :ref:`Environment <environments_resource>`. The list below describes
the questions you will be asked.

.. note::
   If you are not using Vault's AWS secret backend, you can authenticate with AWS in the same way as the AWS CLI. More details can be found at `Configuring the AWS CLI <http://docs.aws.amazon.com /cli/latest/userguide/cli-chap-getting-started.html>`_.

* Configuring a new :ref:`Provider <providers_resource>`
   * Provider name: must be unique
   * Cloud: Amazon (AWS) is the default and only option for now (more clouds to come)
   * Credentials: Amazon CLI auth (i.e. env variables/profile) or Vault (optional)
   * Name prefix: for state buckets and DynamoDB tables
   * Public DNS zone: will be created if not already existing, must be delegated from the root

* Configuring a new :ref:`Environment <environments_resource>`
   * Environment name: must be unique
   * Project name: used for AWS resource labels
   * Project administrator mail address
   * Cloud region: pick a region fetched from AWS (using Provider credentials)

* Configuring new :ref:`Cluster(s) <clusters_resource>`
   * Single or multi-cluster environment
   * Cloud availability zone(s): pick zone(s) fetched from AWS

Once initialised, the configuration will be created at ``$HOME/.tarmak/tarmak.yaml`` (default).

.. _create_ami:

Create an AMI
~~~~~~~~~~~~~
Next we create an AMI for this environment by running ``tarmak clusters images build`` (this is the step that requires Docker to be installed locally).

::

  % tarmak clusters images build
  <output omitted>

.. _create_cluster:

Create the cluster
~~~~~~~~~~~~~~~~~~
To create the cluster, run ``tarmak clusters apply``.

::

  % tarmak clusters apply
  <output omitted>

.. warning::
   The first time this command is run, Tarmak will create a `hosted zone <http://docs.aws.amazon.com/Route53/latest/DeveloperGuide/CreatingHostedZone.html>`_ and then fail with the following error.

   ::

      * failed verifying delegation of public zone 5 times, make sure the zone k8s.jetstack.io is delegated to nameservers [ns-100.awsdns-12.com ns-1283.awsdns-32.org ns-1638.awsdns-12.co.uk ns-842.awsdns-41.net]

You should now change the nameservers of your domain to the four listed in the error. If you only wish to delegate a subdomain containing your zone to AWS without delegating the parent domain see `Creating a Subdomain That Uses Amazon Route 53 as the DNS Service without Migrating the Parent Domain <http://docs.aws.amazon.com/Route53/latest/DeveloperGuide/CreatingNewSubdomain.html>`_.

To complete the cluster provisioning, run ``tarmak clusters apply`` once again.

.. note::
   This process may take 30-60 minutes to complete.
   You can stop it by sending the signal `SIGTERM` or `SIGINT` (Ctrl-C) to the process.
   Tarmak will not exit immediately.
   It will wait for the currently running step to finish and then exit.
   You can complete the process by re-running the command.

.. _use_cluster:

Use the cluster
~~~~~~~~~~~~~~~

Now it is time to take your new cluster for a test drive.
We'll deploy `Sock Shop`_, a sample microservices application that shows how to run and connect a set of services on Kubernetes.

Tarmak provides a convenient ``kubectl`` wrapper which automatically connects to the current Tarmak cluster.
We'll use this instead of running ``kubectl`` directly.

Create a namespace and then deploy the application::

  % tarmak clusters kubectl create namespace sock-shop
  <output omitted>

  % tarmak clusters kubectl apply -n sock-shop -f "https://github.com/microservices-demo/microservices-demo/blob/master/deploy/kubernetes/complete-demo.yaml?raw=true"
  <output omitted>

It takes several minutes to download and start all the containers, watch the output of ``tarmak clusters kubectl get pods -n sock-shop`` to see when they’re all up and running.

You can then find out the port that is allocated for the Sock Shop front-end service by running::

  % tarmak clusters kubectl -n sock-shop get svc front-end

  NAME        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
  front-end   10.254.176.100   <nodes>       80:30001/TCP   21m

  You should be able to connect to the Sock Shop front end on port 30001 of the ingress IP that Tarmak has set up.

  Visit http://sockshop.cluster.richardw.dev.tarmak.org:30001/ in your web browser.

OR

Create an SSH tunnel to the Sock Shop node port on the Kubernetes master node::
  % tarmak clusters ssh -L30001:master:30001 bastion

Then visit http://localhost:30001 in your web browser.

.. _destroy_cluster:

Destroy the cluster
~~~~~~~~~~~~~~~~~~~
To destroy the cluster, run ``tarmak clusters destroy``.

::

  % tarmak clusters destroy
  <output omitted>

.. note::
   This process may take 30-60 minutes to complete.
   You can stop it by sending the signal ``SIGTERM`` or ``SIGINT`` (Ctrl-C) to the process.
   Tarmak will not exit immediately.
   It will wait for the currently running step to finish and then exit.
   You can complete the process by re-running the command.
