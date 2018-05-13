Kube2IAM
--------

Setup terraform IAM policies
****************************

::

  x


Setup kube2iam DaemonSet
************************

::

  helm upgrade kube2iam stable/kube2iam \
    --install \
    --version 0.8.7 \
    --namespace kube-system \
    --set=extraArgs.base-role-arn=arn:aws:iam::0123456789:role/ \
    --set=extraArgs.host-ip=127.0.0.1 \
    --set=extraArgs.log-format=json \
    --set=updateStrategy=RollingUpdate \
    --set=rbac.create=true \
    --set=host.iptables=false

