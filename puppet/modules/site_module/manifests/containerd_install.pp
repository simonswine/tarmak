# Class install containerd
class site_module::containerd_install (
  String $download_url = 'https://storage.googleapis.com/cri-containerd-release/cri-containerd-#VERSION#.linux-amd64.tar.gz',
  String $cni_download_url = 'https://storage.googleapis.com/cri-containerd-release/cri-containerd-cni-#VERSION#.linux-amd64.tar.gz',
){
  include site_module::containerd

  $version  = $::site_module::containerd::version
  $path     = $::site_module::containerd::path
  $dest_dir = $::site_module::containerd::dest_dir

  $_download_url = regsubst(
    $download_url,
    '#VERSION#',
    $version,
    'G'
  )
  $_cni_download_url = regsubst(
    $cni_download_url,
    '#VERSION#',
    $version,
    'G'
  )

  ensure_resource('package', ['curl'],{
    ensure => present
  })

  file { $dest_dir:
    ensure => directory,
    mode   => '0755',
  }
  -> file { "${dest_dir}/bin":
    ensure => directory,
    mode   => '0755',
  }
  -> file { "${dest_dir}/cni":
    ensure => directory,
    mode   => '0755',
  }

  $containerd_tar_path = "${dest_dir}/containerd.tar.gz"
  $containerd_cni_tar_path = "${dest_dir}/containerd-cni.tar.gz"

  File[$dest_dir]
  -> exec {"containerd-${version}-download":
    command => "curl -sL -o  ${containerd_tar_path} ${_download_url}",
    creates => $containerd_tar_path,
    path    => $path,
    require => Package['curl'],
  } -> exec {"containerd-${version}-extract":
    command => "tar xvfz ${containerd_tar_path} -C ${dest_dir}/bin --strip-components=4 ./usr/local/bin ./usr/local/sbin",
    path    => $path,
    require => [
      File["${dest_dir}/bin"],
      Package['curl'],
    ],
    creates => "${dest_dir}/bin/containerd",
  }

  File[$dest_dir]
  -> exec {"containerd-cni-${version}-download":
    command => "curl -sL -o  ${containerd_cni_tar_path} ${_cni_download_url}",
    creates => $containerd_cni_tar_path,
    path    => $path,
    require => Package['curl'],
  } -> exec {"containerd-cni-${version}-extract":
    command => "tar xvfz ${containerd_cni_tar_path} -C ${dest_dir}/cni --strip-components=3 ./opt/cni",
    path    => $path,
    require => [
      File["${dest_dir}/cni"],
      Package['curl'],
    ],
    creates => "${dest_dir}/cni/bin/host-local",
  }

  $bin_dir = '/opt/bin'
  ensure_resource('file', $bin_dir, {
    ensure => directory,
    mode   => '0755',
  })

  File[$bin_dir]
  -> file { "${bin_dir}/crictl":
    ensure => link,
    target => "${dest_dir}/bin/crictl",
  }

  File[$bin_dir]
  -> file { "${bin_dir}/ctr":
    ensure => link,
    target => "${dest_dir}/bin/ctr",
  }

  file { '/opt/cni':
    ensure => link,
    target => "${dest_dir}/cni",
  }

}
