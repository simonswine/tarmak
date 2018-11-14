local kubernetes = import "kubernetes-mixin/mixin.libsonnet";

kubernetes {
  _config+:: {
    kubeStateMetricsSelector: 'job="kubernetes-service-endpoints",app="kube-state-metrics"',
    cadvisorSelector: 'job="kubernetes-nodes-cadvisor"',
    nodeExporterSelector: 'job=~"(bastion|vault|etcd|kubernetes)-nodes-exporter"',
    kubeletSelector: 'job="kubernetes-nodes-cadvisor"',
    kubeApiserverSelector: 'job="kubernetes-apiservers"',
    kubeControllerManagerSelector: 'job="kubernetes-controller-managers"',
    kubeSchedulerSelector: 'job="kubernetes-schedulers"',
  },
}
