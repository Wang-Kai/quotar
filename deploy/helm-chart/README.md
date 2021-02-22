# nfs-provisioner helm chart

## install

```bash
helm install nfs-quotar -n MY_NAMESPACE .
```

## create PVC with nfs-quotar storage class

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: quotar-pvc-test
  annotations:
    volume.beta.kubernetes.io/storage-class: "nfs-with-quota"
spec:
  accessModes:
  - ReadWriteMany
  resources:
    requests:
      storage: 111m
```

```bash
kubectl create -f  pvc.yaml
```



