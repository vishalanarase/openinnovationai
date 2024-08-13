# openinnovationai

## Cluster Setup
> Setup the kind cluster by using kind [configuration](docs/kind.yaml)

```bash
❯ kind create cluster --config kind.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.30.2) 🖼
 ✓ Preparing nodes 📦 📦 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind
```

> Verify the cluster
```bash
❯ k get no
NAME                 STATUS   ROLES           AGE     VERSION
kind-control-plane   Ready    control-plane   5m6s    v1.30.2
kind-worker          Ready    <none>          4m47s   v1.30.2
kind-worker2         Ready    <none>          4m47s   v1.30.2
```

## Install CRDs
> Goto the operator directory

`$ cd distributed-job-scheduler-operator`

```bash
❯ make install
test -s /Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen && /Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen --version | grep -q v0.13.0 || \
	GOBIN=/Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/computejobs.infra.openinnovation.ai created
customresourcedefinition.apiextensions.k8s.io/computenodes.infra.openinnovation.ai created
```

## Run operator locally
> Goto the operator directory

`$ cd distributed-job-scheduler-operator`

> Keep operator running on one terminal
```bash
❯ make run
test -s /Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen && /Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen --version | grep -q v0.13.0 || \
	GOBIN=/Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go run ./cmd/main.go
2024-08-13T18:21:27+05:30	INFO	setup	starting manager
2024-08-13T18:21:27+05:30	INFO	starting server	{"kind": "health probe", "addr": "[::]:8081"}
2024-08-13T18:21:27+05:30	INFO	controller-runtime.metrics	Starting metrics server
2024-08-13T18:21:27+05:30	INFO	controller-runtime.metrics	Serving metrics server	{"bindAddress": ":8080", "secure": false}
2024-08-13T18:21:27+05:30	INFO	Starting EventSource	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "source": "kind source: *v1.ComputeNode"}
2024-08-13T18:21:27+05:30	INFO	Starting EventSource	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "source": "kind source: *v1.Node"}
2024-08-13T18:21:27+05:30	INFO	Starting Controller	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode"}
2024-08-13T18:21:27+05:30	INFO	Starting EventSource	{"controller": "computejob", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeJob", "source": "kind source: *v1.ComputeJob"}
2024-08-13T18:21:27+05:30	INFO	Starting EventSource	{"controller": "computejob", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeJob", "source": "kind source: *v1.Pod"}
2024-08-13T18:21:27+05:30	INFO	Starting Controller	{"controller": "computejob", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeJob"}
2024-08-13T18:21:27+05:30	INFO	Starting workers	{"controller": "computejob", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeJob", "worker count": 1}
2024-08-13T18:21:27+05:30	INFO	Starting workers	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "worker count": 1}
2024-08-13T18:21:27+05:30	INFO	Created new compute node	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "ComputeNode": {"name":"kind-control-plane"}, "namespace": "", "name": "kind-control-plane", "reconcileID": "d4331628-885a-45ff-aae5-d975ec952866", "Name": "kind-control-plane"}
2024-08-13T18:21:27+05:30	INFO	Updated compute node status	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "ComputeNode": {"name":"kind-control-plane"}, "namespace": "", "name": "kind-control-plane", "reconcileID": "d4331628-885a-45ff-aae5-d975ec952866", "Name": "kind-control-plane", "State": "Running"}
2024-08-13T18:21:27+05:30	INFO	Created new compute node	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "ComputeNode": {"name":"kind-control-plane"}, "namespace": "", "name": "kind-control-plane", "reconcileID": "d4331628-885a-45ff-aae5-d975ec952866", "Name": "kind-worker"}
2024-08-13T18:21:27+05:30	INFO	Updated compute node status	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "ComputeNode": {"name":"kind-control-plane"}, "namespace": "", "name": "kind-control-plane", "reconcileID": "d4331628-885a-45ff-aae5-d975ec952866", "Name": "kind-worker", "State": "Running"}
2024-08-13T18:21:27+05:30	INFO	Created new compute node	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "ComputeNode": {"name":"kind-control-plane"}, "namespace": "", "name": "kind-control-plane", "reconcileID": "d4331628-885a-45ff-aae5-d975ec952866", "Name": "kind-worker2"}
2024-08-13T18:21:27+05:30	INFO	Updated compute node status	{"controller": "computenode", "controllerGroup": "infra.openinnovation.ai", "controllerKind": "ComputeNode", "ComputeNode": {"name":"kind-control-plane"}, "namespace": "", "name": "kind-control-plane", "reconcileID": "d4331628-885a-45ff-aae5-d975ec952866", "Name": "kind-worker2", "State": "Running"}
```

> Verify that the computenodes are created
```bash
❯ k get computenodes
NAME                 AGE    STATUS
kind-control-plane   100s   Running
kind-worker          100s   Running
kind-worker2         100s   Running
```

## Create computejob
> On another reminal create resource and monitor states

```yaml
❯ cat config/samples/infra_v1_computejob.yaml
apiVersion: infra.openinnovation.ai/v1
kind: ComputeJob
metadata:
  labels:
    app.kubernetes.io/name: computejob
    app.kubernetes.io/instance: computejob-sample
    app.kubernetes.io/part-of: distributed-job-scheduler-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: distributed-job-scheduler-operator
  name: computejob-sample
spec:
  command: "echo 'Hello, World!' && sleep 20 && exit 0"
  nodeSelector:
    kubernetes.io/os: linux
  parallelism: 2
```

> Create compute job
```bash
❯ k create -f config/samples/infra_v1_computejob.yaml
computejob.infra.openinnovation.ai/computejob-sample created
```

> Verify the compute job
```bash
❯ k get computejobs
NAME                AGE   STATUS
computejob-sample   62s   Running
```

> Verify the pods are created
```bash
❯ k get po
NAME                                   READY   STATUS    RESTARTS   AGE
computejob-sample-kind-control-plane   1/1     Running   0          11s
computejob-sample-kind-worker          1/1     Running   0          11s
```

> Verify compute job is completed
```bash
❯ k get computejobs
NAME                AGE     STATUS
computejob-sample   2m41s   Completed
```

> Verify pods are completed
```bash
❯ k get po
NAME                                   READY   STATUS      RESTARTS   AGE
computejob-sample-kind-control-plane   0/1     Completed   0          111s
computejob-sample-kind-worker          0/1     Completed   0          111s
```

### Operator Setup
- [initalize operator](docs/operator-init.md)
