
generate:
	mkdir -p out
	mkdir -p out/dashboards
	jsonnet -J vendor -S -e 'std.manifestYamlDoc((import "mixin.libsonnet").prometheusAlerts)' > out/alerts.yml
	jsonnet -J vendor -S -e 'std.manifestYamlDoc((import "mixin.libsonnet").prometheusRules)' > out/rules.yml
	jsonnet -J vendor -m out/dashboards -e '(import "mixin.libsonnet").grafanaDashboards'

