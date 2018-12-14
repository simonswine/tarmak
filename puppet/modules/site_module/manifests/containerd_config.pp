# Class install containerd::config
class site_module::containerd_config {
  include site_module::containerd

  $dest_dir = $::site_module::containerd::dest_dir

  file { '/etc/crictl.yaml':
    ensure  => file,
    content => template('site_module/crictl.yaml.erb'),
  }

  file { '/etc/containerd':
    ensure => directory,
    mode   => '0755',
  }
  -> file { '/etc/containerd/config.toml':
    ensure  => file,
    content => template('site_module/containerd_config.toml.erb'),
  }

  file { "${dest_dir}/cni.template":
    ensure  => file,
    content => template('site_module/cni.template.erb'),
  }

}
