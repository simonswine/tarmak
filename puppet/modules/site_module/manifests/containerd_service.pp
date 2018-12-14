class site_module::containerd_service {
  include site_module::containerd

  $dest_dir = $::site_module::containerd::dest_dir
  $service_name = 'containerd'

  exec { "${service_name}-systemctl-daemon-reload":
    command     => '/bin/systemctl daemon-reload',
    refreshonly => true,
    path        => $::site_module::containerd::path
  }

  file { "${::site_module::containerd::systemd_dir}/${service_name}.service":
    ensure  => file,
    content => template('site_module/containerd.service.erb'),
    notify  => Exec["${service_name}-systemctl-daemon-reload"]
  }
  ~> service { "${service_name}.service":
    ensure  => running,
    enable  => true,
    require => Exec["${service_name}-systemctl-daemon-reload"],
  }
}
