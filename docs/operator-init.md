# Init operator repo
```bash
❯ operator-sdk init --domain openinnovation.ai --repo github.com/vishalanarase/openinnovationai/distributed-job-scheduler-operator
INFO[0000] Writing kustomize manifests for you to edit...
INFO[0000] Writing scaffold for you to edit...
INFO[0000] Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.16.3
INFO[0002] Update dependencies:
$ go mod tidy
Next: define a resource with:
$ operator-sdk create api
```

## Create ComputeNode API, Resource, Controller
```bash
❯ operator-sdk create api --group infra --version v1 --kind ComputeNode --resource --controller
INFO[0000] Writing kustomize manifests for you to edit...
INFO[0000] Writing scaffold for you to edit...
INFO[0000] api/v1/computenode_types.go
INFO[0000] api/v1/groupversion_info.go
INFO[0000] internal/controller/suite_test.go
INFO[0000] internal/controller/computenode_controller.go
INFO[0000] internal/controller/computenode_controller_test.go
INFO[0000] Update dependencies:
$ go mod tidy
INFO[0000] Running make:
$ make generate
/Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

## Create ComputeJob API, Resource, Controller
```bash
❯ operator-sdk create api --group infra --version v1 --kind ComputeJob --resource --controller
INFO[0000] Writing kustomize manifests for you to edit...
INFO[0000] Writing scaffold for you to edit...
INFO[0000] api/v1/computejob_types.go
INFO[0000] api/v1/groupversion_info.go
INFO[0000] internal/controller/suite_test.go
INFO[0000] internal/controller/computejob_controller.go
INFO[0000] internal/controller/computejob_controller_test.go
INFO[0000] Update dependencies:
$ go mod tidy
INFO[0000] Running make:
$ make generate
/Users/vishal/workspace/vishalanarase/openinnovationai/distributed-job-scheduler-operator/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```