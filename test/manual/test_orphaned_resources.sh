# -------- SETUP --------

# 1. Create a test namespace
kubectl create namespace test-orphanage

# 2. Create an OrphanagePolicy to monitor Secrets and ConfigMaps
kubectl apply -f - <<EOF
apiVersion: orphanage.kponos.io/v1alpha1
kind: OrphanagePolicy
metadata:
  name: test-policy
  namespace: test-orphanage
spec:
  resourceTypes:
    - Secret
    - ConfigMap
EOF

# 3. Check if the policy was created
kubectl get orphanagepolicy -n test-orphanage

# 4. Check the status of the policy (should show orphans)
kubectl get orphanagepolicy test-policy -n test-orphanage -o yaml
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[*].name}'



# -------- CASE 1: Orphaned Secret --------

# Create an orphaned secret (not referenced by any resource)
kubectl create secret generic orphaned-secret \
  --from-literal=key1=value1 \
  -n test-orphanage

# Wait a moment for reconciliation, then check status
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[?(@.name=="orphaned-secret")]}'



# -------- CASE 2: Orphaned ConfigMap --------

# Create an orphaned configmap
kubectl create configmap orphaned-configmap \
  --from-literal=key1=value1 \
  -n test-orphanage

# Check if it's detected
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[?(@.name=="orphaned-configmap")]}'



# -------- CASE 3: Secret referenced by Deployment --------

# Create a secret
kubectl create secret generic used-secret \
  --from-literal=password=secret123 \
  -n test-orphanage

# Create a deployment that uses the secret
kubectl create deployment test-deployment \
  --image=nginx:latest \
  -n test-orphanage

# Use secret as volume
kubectl patch deployment test-deployment -n test-orphanage -p '{"spec":{"template":{"spec":{"volumes":[{"name":"secret-vol","secret":{"secretName":"used-secret"}}]}}}}'

# Verify the secret is NOT in orphans list
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[?(@.name=="used-secret")]}'
# Should return nothing



# -------- CASE 4: ConfigMap referenced by Pod --------

# Create a configmap
kubectl create configmap used-configmap \
  --from-literal=app.properties=key=value \
  -n test-orphanage

# Create a pod that uses the configmap
kubectl run test-pod \
  --image=nginx:latest \
  --restart=Never \
  -n test-orphanage \
  --overrides='{"spec":{"volumes":[{"name":"config-vol","configMap":{"name":"used-configmap"}}],"containers":[{"name":"test-pod","image":"nginx:latest","volumeMounts":[{"name":"config-vol","mountPath":"/etc/config"}]}]}}'

# Verify it's NOT orphaned
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[?(@.name=="used-configmap")]}'



# -------- CASE 5: Secret referenced by ServiceAccount --------

# Create a secret for service account
kubectl create secret generic sa-secret \
  --from-literal=token=abc123 \
  -n test-orphanage

# Create a service account that uses the secret
kubectl create serviceaccount test-sa -n test-orphanage
kubectl patch serviceaccount test-sa -n test-orphanage -p '{"secrets":[{"name":"sa-secret"}]}'

# Verify it's NOT orphaned
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[?(@.name=="sa-secret")]}'



# -------- CASE 6: ConfigMap referenced by StatefulSet --------

# Create configmap
kubectl create configmap statefulset-config \
  --from-literal=config.yaml="key: value" \
  -n test-orphanage

# Create a statefulset using the configmap
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-statefulset
  namespace: test-orphanage
spec:
  serviceName: test-statefulset
  replicas: 1
  selector:
    matchLabels:
      app: test-statefulset
  template:
    metadata:
      labels:
        app: test-statefulset
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        volumeMounts:
        - name: config
          mountPath: /etc/config
      volumes:
      - name: config
        configMap:
          name: statefulset-config
EOF

# Verify it's NOT orphaned
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[?(@.name=="statefulset-config")]}'



# -------- CASE 7: Secret referenced by Ingress --------

# Create TLS secret directly using YAML
kubectl create secret generic ingress-tls-secret \
  --from-literal=tls.crt="dummy-cert" \
  --from-literal=tls.key="dummy-key" \
  -n test-orphanage

# Create ingress using the secret
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress
  namespace: test-orphanage
spec:
  tls:
  - hosts:
    - example.com
    secretName: ingress-tls-secret
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-service
            port:
              number: 80
EOF

# Verify it's NOT orphaned
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[?(@.name=="ingress-tls-secret")]}'



# -------- VERIFICATION --------

# Watch the policy status in real-time
kubectl get orphanagepolicy test-policy -n test-orphanage -w

# Get detailed status with all orphans
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status}' | jq

# Count total orphans
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphanCount}'

# List all orphan names
kubectl get orphanagepolicy test-policy -n test-orphanage -o jsonpath='{.status.orphans[*].name}'

# Check operator logs to see reconciliation
kubectl logs -n kponos-system -l control-plane=controller-manager --tail=50

# Verify operator is running
kubectl get pods -n kponos-system



# -------- CLEANUP --------

# Delete test resources
kubectl delete namespace test-orphanage

# Or delete specific resources
kubectl delete orphanagepolicy test-policy -n test-orphanage
kubectl delete secret orphaned-secret -n test-orphanage
kubectl delete configmap orphaned-configmap -n test-orphanage