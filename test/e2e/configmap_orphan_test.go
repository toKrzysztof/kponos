package e2e

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/toKrzysztof/kponos/test/utils"
)

var _ = Describe("ConfigMap Orphan Detection", Ordered, func() {
	var testNamespace string

	BeforeAll(func() {
		// Create a test namespace
		testNamespace = "configmap-orphan-test-" + utils.RandomString(5)
		cmd := exec.Command("kubectl", "create", "ns", testNamespace)
		_, err := utils.Run(cmd)

		Expect(err).NotTo(HaveOccurred(), "Failed to create test namespace")

		// Create OrphanagePolicy for ConfigMaps
		policyYAML := fmt.Sprintf(`apiVersion: orphanage.kponos.io/v1alpha1
kind: OrphanagePolicy
metadata:
  name: configmap-test-policy
  namespace: %s
spec:
  resourceTypes:
    - ConfigMap`, testNamespace)

		cmd = exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = strings.NewReader(policyYAML)
		_, err = utils.Run(cmd)

		Expect(err).NotTo(HaveOccurred(), "Failed to create OrphanagePolicy")
	})

	AfterAll(func() {
		// Clean up test namespace
		cmd := exec.Command("kubectl", "delete", "ns", testNamespace, "--ignore-not-found=true")
		_, _ = utils.Run(cmd)
	})

	SetDefaultEventuallyTimeout(2 * time.Minute)
	SetDefaultEventuallyPollingInterval(2 * time.Second)

	type testCase struct {
		resourceType  string // Pod, Deployment, StatefulSet, DaemonSet
		configMapName string
		resourceName  string
		referencePath string // description of the reference path
		resourceYAML  func(namespace, resourceName, configMapName string) string
		configMapData []string // --from-literal args
	}

	DescribeTable("should NOT detect ConfigMap as orphaned when referenced",
		func(tc testCase) {
			By(fmt.Sprintf("creating ConfigMap %s", tc.configMapName))
			args := append([]string{"create", "configmap", tc.configMapName, "-n", testNamespace}, tc.configMapData...)
			cmd := exec.Command("kubectl", args...)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By(fmt.Sprintf("creating %s %s that references ConfigMap %s (path: %s)", tc.resourceType, tc.resourceName, tc.configMapName, tc.referencePath))
			resourceYAML := tc.resourceYAML(testNamespace, tc.resourceName, tc.configMapName)
			cmd = exec.Command("kubectl", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(resourceYAML)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By(fmt.Sprintf("verifying ConfigMap %s is NOT in orphans list (tested via %s)", tc.configMapName, tc.referencePath))
			verifyNotOrphaned := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "orphanagepolicy", "configmap-test-policy",
					"-n", testNamespace,
					"-o", fmt.Sprintf("jsonpath={.status.orphans[?(@.name==\"%s\")]}", tc.configMapName))
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(BeEmpty(), "ConfigMap should not be in orphans list")
			}
			Eventually(verifyNotOrphaned).Should(Succeed())
		},
		Entry("Pod via volume", testCase{
			resourceType:  "Pod",
			configMapName: "pod-volume-configmap",
			resourceName:  "test-pod-volume",
			referencePath: "volumes[].configMap.name",
			configMapData: []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, configMapName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  containers:
  - name: test
    image: busybox:latest
    command: ["sleep", "infinity"]
    volumeMounts:
    - name: config-vol
      mountPath: /etc/config
  volumes:
  - name: config-vol
    configMap:
      name: %s`, resourceName, namespace, configMapName)
			},
		}),
		Entry("Pod via envFrom", testCase{
			resourceType:  "Pod",
			configMapName: "pod-envfrom-configmap",
			resourceName:  "test-pod-envfrom",
			referencePath: "containers[].envFrom[].configMapRef.name",
			configMapData: []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, configMapName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  containers:
  - name: test
    image: busybox:latest
    command: ["sleep", "infinity"]
    envFrom:
    - configMapRef:
        name: %s`, resourceName, namespace, configMapName)
			},
		}),
		Entry("Pod via env.valueFrom", testCase{
			resourceType:  "Pod",
			configMapName: "pod-env-configmap",
			resourceName:  "test-pod-env",
			referencePath: "containers[].env[].valueFrom.configMapKeyRef.name",
			configMapData: []string{"--from-literal=app.properties=key=value"},
			resourceYAML: func(namespace, resourceName, configMapName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  containers:
  - name: test
    image: busybox:latest
    command: ["sleep", "infinity"]
    env:
    - name: APP_PROPERTIES
      valueFrom:
        configMapKeyRef:
          name: %s
          key: app.properties`, resourceName, namespace, configMapName)
			},
		}),
		Entry("Pod initContainer via envFrom", testCase{
			resourceType:  "Pod",
			configMapName: "pod-init-envfrom-configmap",
			resourceName:  "test-pod-init-envfrom",
			referencePath: "initContainers[].envFrom[].configMapRef.name",
			configMapData: []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, configMapName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  initContainers:
  - name: init
    image: busybox:latest
    command: ["sh", "-c", "echo init done"]
    envFrom:
    - configMapRef:
        name: %s
  containers:
  - name: test
    image: busybox:latest
    command: ["sleep", "infinity"]`, resourceName, namespace, configMapName)
			},
		}),
		Entry("Pod initContainer via env.valueFrom", testCase{
			resourceType:  "Pod",
			configMapName: "pod-init-env-configmap",
			resourceName:  "test-pod-init-env",
			referencePath: "initContainers[].env[].valueFrom.configMapKeyRef.name",
			configMapData: []string{"--from-literal=init.properties=key=value"},
			resourceYAML: func(namespace, resourceName, configMapName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  initContainers:
  - name: init
    image: busybox:latest
    command: ["sh", "-c", "echo init done"]
    env:
    - name: INIT_PROPERTIES
      valueFrom:
        configMapKeyRef:
          name: %s
          key: init.properties
  containers:
  - name: test
    image: busybox:latest
    command: ["sleep", "infinity"]`, resourceName, namespace, configMapName)
			},
		}),
		Entry("Deployment via volume", testCase{
			resourceType:  "Deployment",
			configMapName: "deployment-volume-configmap",
			resourceName:  "test-deployment-volume",
			referencePath: "volumes[].configMap.name",
			configMapData: []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, configMapName string) string {
				return fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
spec:
  replicas: 1
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: test
        image: busybox:latest
        command: ["sleep", "infinity"]
        volumeMounts:
        - name: config-vol
          mountPath: /etc/config
      volumes:
      - name: config-vol
        configMap:
          name: %s`, resourceName, namespace, resourceName, resourceName, configMapName)
			},
		}),
		Entry("StatefulSet via volume", testCase{
			resourceType:  "StatefulSet",
			configMapName: "statefulset-volume-configmap",
			resourceName:  "test-statefulset-volume",
			referencePath: "volumes[].configMap.name",
			configMapData: []string{"--from-literal=config.yaml=key:value"},
			resourceYAML: func(namespace, resourceName, configMapName string) string {
				return fmt.Sprintf(`apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: %s
  namespace: %s
spec:
  serviceName: %s
  replicas: 1
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: test
        image: busybox:latest
        command: ["sleep", "infinity"]
        volumeMounts:
        - name: config
          mountPath: /etc/config
      volumes:
      - name: config
        configMap:
          name: %s`, resourceName, namespace, resourceName, resourceName, resourceName, configMapName)
			},
		}),
		Entry("DaemonSet via volume", testCase{
			resourceType:  "DaemonSet",
			configMapName: "daemonset-volume-configmap",
			resourceName:  "test-daemonset-volume",
			referencePath: "volumes[].configMap.name",
			configMapData: []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, configMapName string) string {
				return fmt.Sprintf(`apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: %s
  namespace: %s
spec:
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: test
        image: busybox:latest
        command: ["sleep", "infinity"]
        volumeMounts:
        - name: config-vol
          mountPath: /etc/config
      volumes:
      - name: config-vol
        configMap:
          name: %s`, resourceName, namespace, resourceName, resourceName, configMapName)
			},
		}),
	)
})
