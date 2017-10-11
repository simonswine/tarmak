// Copyright Jetstack Ltd. See LICENSE for details.
package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ligurio/go-junit/parser"
	"github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/jetstack/tarmak/pkg/tarmak"
	"github.com/jetstack/tarmak/pkg/tarmak/kubectl"
)

var bigLumpOfYaml = `---
apiVersion: v1
kind: Namespace
metadata:
  name: heptio-sonobuoy
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-serviceaccount
  namespace: heptio-sonobuoy
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-serviceaccount
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: sonobuoy-serviceaccount
subjects:
- kind: ServiceAccount
  name: sonobuoy-serviceaccount
  namespace: heptio-sonobuoy
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-serviceaccount
  namespace: heptio-sonobuoy
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
---
apiVersion: v1
data:
  config.json: |
    {
        "Description": "EXAMPLE",
        "Filters": {
            "LabelSelector": "",
            "Namespaces": ".*"
        },
        "PluginNamespace": "heptio-sonobuoy",
        "Plugins": [
            {
                "name": "systemd_logs"
            },
            {
                "name": "e2e"
            }
        ],
        "Resources": [
            "CertificateSigningRequests",
            "ClusterRoleBindings",
            "ClusterRoles",
            "ComponentStatuses",
            "CustomResourceDefinitions",
            "Nodes",
            "PersistentVolumes",
            "PodSecurityPolicies",
            "ServerVersion",
            "StorageClasses",
            "ConfigMaps",
            "DaemonSets",
            "Deployments",
            "Endpoints",
            "Events",
            "HorizontalPodAutoscalers",
            "Ingresses",
            "Jobs",
            "LimitRanges",
            "PersistentVolumeClaims",
            "Pods",
            "PodDisruptionBudgets",
            "PodTemplates",
            "ReplicaSets",
            "ReplicationControllers",
            "ResourceQuotas",
            "RoleBindings",
            "Roles",
            "ServerGroups",
            "ServiceAccounts",
            "Services",
            "StatefulSets"
        ],
        "ResultsDir": "/tmp/sonobuoy",
        "Server": {
            "advertiseaddress": "sonobuoy-master:8080",
            "bindaddress": "0.0.0.0",
            "bindport": 8080,
            "timeoutseconds": 21600
        },
        "Version": "v0.9.0"
    }
kind: ConfigMap
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-config-cm
  namespace: heptio-sonobuoy
---
apiVersion: v1
data:
  e2e.yaml: |
    driver: Job
    name: e2e
    resultType: e2e
    spec:
      containers:
      - env:
        - name: E2E_FOCUS
          value: "Conformance"
        image: gcr.io/heptio-images/kube-conformance:v1.8.0
        imagePullPolicy: Always
        name: e2e
        volumeMounts:
        - mountPath: /tmp/results
          name: results
      - command:
        - sh
        - -c
        - /sonobuoy worker global -v 5 --logtostderr
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: spec.nodeName
        - name: RESULTS_DIR
          value: /tmp/results
        image: gcr.io/heptio-images/sonobuoy:master
        imagePullPolicy: Always
        name: sonobuoy-worker
        volumeMounts:
        - mountPath: /etc/sonobuoy
          name: config
        - mountPath: /tmp/results
          name: results
      restartPolicy: Never
      serviceAccountName: sonobuoy-serviceaccount
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: Exists
      - key: CriticalAddonsOnly
        operator: Exists
      volumes:
      - emptyDir: {}
        name: results
      - configMap:
          name: __SONOBUOY_CONFIGMAP__
        name: config
  systemdlogs.yaml: |
    driver: DaemonSet
    name: systemd_logs
    resultType: systemd_logs
    spec:
      containers:
      - command:
        - sh
        - -c
        - /get_systemd_logs.sh && sleep 21600
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: spec.nodeName
        - name: RESULTS_DIR
          value: /tmp/results
        - name: CHROOT_DIR
          value: /node
        image: gcr.io/heptio-images/sonobuoy-plugin-systemd-logs:latest
        imagePullPolicy: Always
        name: systemd-logs
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /node
          name: root
        - mountPath: /tmp/results
          name: results
        - mountPath: /etc/sonobuoy
          name: config
      - command:
        - sh
        - -c
        - /sonobuoy worker single-node -v 5 --logtostderr && sleep 21600
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: spec.nodeName
        - name: RESULTS_DIR
          value: /tmp/results
        image: gcr.io/heptio-images/sonobuoy:master
        imagePullPolicy: Always
        name: sonobuoy-worker
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /tmp/results
          name: results
        - mountPath: /etc/sonobuoy
          name: config
      dnsPolicy: ClusterFirstWithHostNet
      hostIPC: true
      hostNetwork: true
      hostPID: true
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: Exists
      - key: CriticalAddonsOnly
        operator: Exists
      volumes:
      - hostPath:
          path: /
        name: root
      - emptyDir: {}
        name: results
      - configMap:
          name: __SONOBUOY_CONFIGMAP__
        name: config
kind: ConfigMap
metadata:
  labels:
    component: sonobuoy
  name: sonobuoy-plugins-cm
  namespace: heptio-sonobuoy
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    component: sonobuoy
    run: sonobuoy-master
    tier: analysis
  name: sonobuoy
  namespace: heptio-sonobuoy
spec:
  containers:
  - command:
    - /bin/bash
    - -c
    - /sonobuoy master --no-exit=true -v 3 --logtostderr
    env:
    - name: SONOBUOY_ADVERTISE_IP
      valueFrom:
        fieldRef:
          fieldPath: status.podIP
    image: gcr.io/heptio-images/sonobuoy:master
    imagePullPolicy: Always
    name: kube-sonobuoy
    volumeMounts:
    - mountPath: /etc/sonobuoy
      name: sonobuoy-config-volume
    - mountPath: /plugins.d
      name: sonobuoy-plugins-volume
    - mountPath: /tmp/sonobuoy
      name: output-volume
  restartPolicy: Never
  serviceAccountName: sonobuoy-serviceaccount
  volumes:
  - configMap:
      name: sonobuoy-config-cm
    name: sonobuoy-config-volume
  - configMap:
      name: sonobuoy-plugins-cm
    name: sonobuoy-plugins-volume
  - emptyDir: {}
    name: output-volume
---
apiVersion: v1
kind: Service
metadata:
  labels:
    component: sonobuoy
    run: sonobuoy-master
  name: sonobuoy-master
  namespace: heptio-sonobuoy
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    run: sonobuoy-master
  type: ClusterIP
`

func getE2eLogsStream(client *kubernetes.Clientset, pod *v1.Pod) (io.ReadCloser, error) {
	logOptions := &v1.PodLogOptions{
		Container:  pod.Spec.Containers[0].Name,
		Follow:     false,
		Previous:   false,
		Timestamps: true,
	}

	return client.CoreV1().RESTClient().Get().
		Namespace(pod.Namespace).
		Name(pod.Name).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).Stream()
}

func isSonobuoyComplete(client *kubernetes.Clientset, pod *v1.Pod) (bool, error) {
	logStream, err := getE2eLogsStream(client, pod)
	if err != nil {
		return false, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(logStream)
	logStr := buf.String()

	return strings.Contains(logStr, "msg=\"no-exit was specified, sonobuoy is now blocking\""), nil
}

func parseReport(r io.Reader) (*junit.JUnitTestsuite, error) {
	var report = new(junit.JUnitTestsuite)

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal([]byte(buf), &report)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func getJunitFromTarball(gzipStream io.Reader) (*junit.JUnitTestsuite, error) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if err != nil {
				return nil, err
			}
			if header.Name == "plugins/e2e/results/junit_01.xml" {
				return parseReport(tarReader)
			}
		default:
			return nil, errors.New("unknown type in tarball")
		}
	}
	return nil, errors.New("e2e test results not found")
}

func findSonobouyTarball(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "sonobuoy") && strings.HasSuffix(f.Name(), "tar.gz") {
			return path.Join(dir, f.Name()), nil
		}
	}
	return "", errors.New("sonobuoy tarball not found")
}

func getE2EResults() (*junit.JUnitTestsuite, error) {
	filename, err := findSonobouyTarball("/tmp/sonobuoy/")
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return getJunitFromTarball(file)
}

func showE2ESummary(t *junit.JUnitTestsuite) {
	if t.Failures > 0 {
		fmt.Printf("\nRan %d tests, %d failures:\n", t.Tests, t.Failures)

		for _, test := range t.TestCases {
			if test.Failure != nil {
				fmt.Println("  " + test.Name)
				fail := test.Failure
				if fail != nil {
					value := html.UnescapeString(fail.Value)
					if false {
						fmt.Println(value)
					}
				}
			}
		}
	} else {
		fmt.Printf("\nRan %d tests, all successful!\n", t.Tests)
	}
}

var clusterVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run tests to verify the current cluster is fully operational",
	Run: func(cmd *cobra.Command, args []string) {
		t := tarmak.New(globalFlags)
		k := t.Kubectl().(*kubectl.Kubectl)
		k.EnsureConfig()

		err := ioutil.WriteFile("/tmp/sonobuoy.yaml", []byte(bigLumpOfYaml), 0644)
		k.Kubectl([]string{"apply", "--validate=false", "-f", "/tmp/sonobuoy.yaml"})

		config, err := clientcmd.BuildConfigFromFlags("", k.ConfigPath())
		if err != nil {
			fmt.Println(err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
		}

		pods, err := clientset.CoreV1().Pods("heptio-sonobuoy").List(metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
		}
		for _, pod := range pods.Items {
			fmt.Println(pod.Name)
		}

		sonobuoy, err := clientset.CoreV1().Pods("heptio-sonobuoy").Get("sonobuoy", metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("waiting for test results")
		for {
			done, err := isSonobuoyComplete(clientset, sonobuoy)
			if err != nil {
				fmt.Println(err)
			}
			if done {
				break
			} else {
				fmt.Printf(".")
				time.Sleep(5 * time.Second)
			}
		}

		k.Kubectl([]string{"cp", "heptio-sonobuoy/sonobuoy:/tmp/sonobuoy", "/tmp/sonobuoy", "--namespace=heptio-sonobuoy"})
		defer os.Remove("/tmp/sonobuoy")

		testResults, err := getE2EResults()

		k.Kubectl([]string{"delete", "-f", "/tmp/sonobuoy.yaml"})
		os.Remove("/tmp/sonobuoy.yaml")

		showE2ESummary(testResults)
	},
}

func init() {
	clusterCmd.AddCommand(clusterVerifyCmd)
}
