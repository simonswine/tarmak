# Class install site_module::containerd
class site_module::containerd (
  String $version = '1.2.1',
){

  $path = defined('$::path') ? {
    default => '/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/bin',
    true    => $::path
  }

  $dest_dir = "/opt/containerd-${version}"
  $systemd_dir = '/etc/systemd/system'

  class { '::site_module::containerd_install': }
  -> class { '::site_module::containerd_config': }
  ~> class { '::site_module::containerd_service': }
  -> Class['::site_module::containerd']

}
