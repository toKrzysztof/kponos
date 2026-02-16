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

var _ = Describe("Secret Orphan Detection", Ordered, func() {
	var testNamespace string

	BeforeAll(func() {
		// Verify CRD is available (should already be from BeforeSuite, but check as safety net)
		By("verifying OrphanagePolicy CRD is available")
		verifyCRDAvailable := func(g Gomega) {
			cmd := exec.Command("kubectl", "get", "crd", "orphanagepolicies.orphanage.kponos.io")
			_, err := utils.Run(cmd)
			g.Expect(err).NotTo(HaveOccurred(), "OrphanagePolicy CRD not available")
		}
		Eventually(verifyCRDAvailable).WithTimeout(10 * time.Second).WithPolling(500 * time.Millisecond).Should(Succeed())

		// Create a test namespace
		testNamespace = "secret-orphan-test-" + utils.RandomString(5)
		cmd := exec.Command("kubectl", "create", "ns", testNamespace)
		_, err := utils.Run(cmd)

		Expect(err).NotTo(HaveOccurred(), "Failed to create test namespace")

		// Create OrphanagePolicy for Secrets
		policyYAML := fmt.Sprintf(`apiVersion: orphanage.kponos.io/v1alpha1
kind: OrphanagePolicy
metadata:
  name: secret-test-policy
  namespace: %s
spec:
  resourceTypes:
    - Secret`, testNamespace)

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
		secretName    string
		resourceName  string
		referencePath string // description of the reference path
		resourceYAML  func(namespace, resourceName, secretName string) string
		secretData    []string // --from-literal args
	}

	DescribeTable("should NOT detect Secret as orphaned when referenced",
		func(tc testCase) {
			By(fmt.Sprintf("creating Secret %s", tc.secretName))
			args := append([]string{"create", "secret", "generic", tc.secretName, "-n", testNamespace}, tc.secretData...)
			cmd := exec.Command("kubectl", args...)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By(fmt.Sprintf("creating %s %s that references Secret %s (path: %s)", tc.resourceType, tc.resourceName,
				tc.secretName, tc.referencePath))
			resourceYAML := tc.resourceYAML(testNamespace, tc.resourceName, tc.secretName)
			cmd = exec.Command("kubectl", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(resourceYAML)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			By(fmt.Sprintf("verifying Secret %s is NOT in orphans list (tested via %s)", tc.secretName, tc.referencePath))
			verifyNotOrphaned := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "orphanagepolicy", "secret-test-policy",
					"-n", testNamespace,
					"-o", fmt.Sprintf("jsonpath={.status.orphans[?(@.name==\"%s\")]}", tc.secretName))
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(BeEmpty(), "Secret should not be in orphans list")
			}
			Eventually(verifyNotOrphaned).Should(Succeed())
		},
		Entry("Pod via volume", testCase{
			resourceType:  "Pod",
			secretName:    "pod-volume-secret",
			resourceName:  "test-pod-volume",
			referencePath: "volumes[].secret.secretName",
			secretData:    []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
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
    - name: secret-vol
      mountPath: /etc/secret
  volumes:
  - name: secret-vol
    secret:
      secretName: %s`, resourceName, namespace, secretName)
			},
		}),
		Entry("Pod via envFrom", testCase{
			resourceType:  "Pod",
			secretName:    "pod-envfrom-secret",
			resourceName:  "test-pod-envfrom",
			referencePath: "containers[].envFrom[].secretRef.name",
			secretData:    []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
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
    - secretRef:
        name: %s`, resourceName, namespace, secretName)
			},
		}),
		Entry("Pod via env.valueFrom", testCase{
			resourceType:  "Pod",
			secretName:    "pod-env-secret",
			resourceName:  "test-pod-env",
			referencePath: "containers[].env[].valueFrom.secretKeyRef.name",
			secretData:    []string{"--from-literal=app.properties=key=value"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
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
        secretKeyRef:
          name: %s
          key: app.properties`, resourceName, namespace, secretName)
			},
		}),
		Entry("Pod initContainer via envFrom", testCase{
			resourceType:  "Pod",
			secretName:    "pod-init-envfrom-secret",
			resourceName:  "test-pod-init-envfrom",
			referencePath: "initContainers[].envFrom[].secretRef.name",
			secretData:    []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
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
    - secretRef:
        name: %s
  containers:
  - name: test
    image: busybox:latest
    command: ["sleep", "infinity"]`, resourceName, namespace, secretName)
			},
		}),
		Entry("Pod initContainer via env.valueFrom", testCase{
			resourceType:  "Pod",
			secretName:    "pod-init-env-secret",
			resourceName:  "test-pod-init-env",
			referencePath: "initContainers[].env[].valueFrom.secretKeyRef.name",
			secretData:    []string{"--from-literal=init.properties=key=value"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
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
        secretKeyRef:
          name: %s
          key: init.properties
  containers:
  - name: test
    image: busybox:latest
    command: ["sleep", "infinity"]`, resourceName, namespace, secretName)
			},
		}),
		Entry("Pod via imagePullSecrets", testCase{
			resourceType:  "Pod",
			secretName:    "pod-imagepull-secret",
			resourceName:  "test-pod-imagepull",
			referencePath: "imagePullSecrets[].name",
			secretData:    []string{"--from-literal=.dockerconfigjson={\"auths\":{}}"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  imagePullSecrets:
  - name: %s
  containers:
  - name: test
    image: busybox:latest
    command: ["sleep", "infinity"]`, resourceName, namespace, secretName)
			},
		}),
		Entry("Deployment via volume", testCase{
			resourceType:  "Deployment",
			secretName:    "deployment-volume-secret",
			resourceName:  "test-deployment-volume",
			referencePath: "volumes[].secret.secretName",
			secretData:    []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
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
        - name: secret-vol
          mountPath: /etc/secret
      volumes:
      - name: secret-vol
        secret:
          secretName: %s`, resourceName, namespace, resourceName, resourceName, secretName)
			},
		}),
		Entry("StatefulSet via volume", testCase{
			resourceType:  "StatefulSet",
			secretName:    "statefulset-volume-secret",
			resourceName:  "test-statefulset-volume",
			referencePath: "volumes[].secret.secretName",
			secretData:    []string{"--from-literal=config.yaml=key:value"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
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
        - name: secret
          mountPath: /etc/secret
      volumes:
      - name: secret
        secret:
          secretName: %s`, resourceName, namespace, resourceName, resourceName, resourceName, secretName)
			},
		}),
		Entry("DaemonSet via volume", testCase{
			resourceType:  "DaemonSet",
			secretName:    "daemonset-volume-secret",
			resourceName:  "test-daemonset-volume",
			referencePath: "volumes[].secret.secretName",
			secretData:    []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
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
        - name: secret-vol
          mountPath: /etc/secret
      volumes:
      - name: secret-vol
        secret:
          secretName: %s`, resourceName, namespace, resourceName, resourceName, secretName)
			},
		}),
		Entry("Ingress via TLS", testCase{
			resourceType:  "Ingress",
			secretName:    "ingress-tls-secret",
			resourceName:  "test-ingress-tls",
			referencePath: "spec.tls[].secretName",
			secretData:    []string{"--from-literal=tls.crt=cert", "--from-literal=tls.key=key"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  name: test-ingress-backend
  namespace: %s
spec:
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: test
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: %s
  namespace: %s
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - example.com
    secretName: %s
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-ingress-backend
            port:
              number: 80`, namespace, resourceName, namespace, secretName)
			},
		}),
		Entry("ServiceAccount via secrets", testCase{
			resourceType:  "ServiceAccount",
			secretName:    "sa-secret",
			resourceName:  "test-sa-secrets",
			referencePath: "secrets[].name",
			secretData:    []string{"--from-literal=key1=value1"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: ServiceAccount
metadata:
  name: %s
  namespace: %s
secrets:
- name: %s`, resourceName, namespace, secretName)
			},
		}),
		Entry("ServiceAccount via imagePullSecrets", testCase{
			resourceType:  "ServiceAccount",
			secretName:    "sa-imagepull-secret",
			resourceName:  "test-sa-imagepull",
			referencePath: "imagePullSecrets[].name",
			secretData:    []string{"--from-literal=.dockerconfigjson={\"auths\":{}}"},
			resourceYAML: func(namespace, resourceName, secretName string) string {
				return fmt.Sprintf(`apiVersion: v1
kind: ServiceAccount
metadata:
  name: %s
  namespace: %s
imagePullSecrets:
- name: %s`, resourceName, namespace, secretName)
			},
		}),
	)
})
