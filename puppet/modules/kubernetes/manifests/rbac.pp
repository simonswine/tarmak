# This class manages RBAC manifests
class kubernetes::rbac{
  require ::kubernetes

  $authorization_mode = $kubernetes::_authorization_mode
  if member($authorization_mode, 'RBAC') and versioncmp($::kubernetes::version, '1.6.0') < 0 {
    kubernetes::apply{'puppernetes-rbac':
      manifests => [
        template('kubernetes/rbac-namespace-kube-public.yaml.erb'),
        template('kubernetes/rbac-cluster-roles.yaml.erb'),
        template('kubernetes/rbac-cluster-role-bindings.yaml.erb'),
        template('kubernetes/rbac-namespace-roles.yaml.erb'),
        template('kubernetes/rbac-namespace-role-bindings.yaml.erb'),
      ],
    }
  }

  # This is hotfix for enabling k8s 1.8 without node authorizer
  $authorization_mode = $kubernetes::_authorization_mode
  if member($authorization_mode, 'RBAC') and versioncmp($::kubernetes::version, '1.8.0') < 0 {
    kubernetes::apply{'puppernetes-rbac':
      manifests => [
        template('kubernetes/rbac-node-cluster-role-bindings.yaml.erb'),
      ],
    }
  }
}
