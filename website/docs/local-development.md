# Local development quickstart

We will:

* Run ORC locally:
    * Create a local kind cluster
    * Load the ORC CRDs
    * Run the ORC manager locally directly from source
* Create an example ORC resource:
    * Create OpenStack credentials
    * Initialise the kustomize environment and load OpenStack credentials
    * Create an ORC resource using the above

## Run ORC locally

### Create a kind cluster

Obtain `kind` from https://kind.sigs.k8s.io/.

Create a default kind cluster:
```bash
$ kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.30.0) 🖼c
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
```

Ensure your local context is set to the newly created kind cluster:
```bash
$ kubectl get node
NAME                 STATUS   ROLES           AGE     VERSION
kind-control-plane   Ready    control-plane   4m22s   v1.30.0
```

### Load the ORC CRDs

From the root directory:
```bash
$ kubectl apply -k config/crd --server-side
customresourcedefinition.apiextensions.k8s.io/images.openstack.k-orc.cloud serverside-applied
customresourcedefinition.apiextensions.k8s.io/networks.openstack.k-orc.cloud serverside-applied
customresourcedefinition.apiextensions.k8s.io/subnets.openstack.k-orc.cloud serverside-applied
```

### Run the ORC manager locally

```bash
$ go run ./cmd/manager -zap-log-level 5
2024-11-11T12:09:30Z    INFO    setup   starting manager
2024-11-11T12:09:30Z    INFO    starting server {"name": "health probe", "addr": "[::]:8081"}
2024-11-11T12:09:30Z    INFO    Starting EventSource    {"controller": "network", "controllerGroup": "openstack.k-orc.cloud", "controllerKind": "Network", "source": "kind source: *v1alpha1.Network"}
2024-11-11T12:09:30Z    INFO    Starting Controller     {"controller": "network", "controllerGroup": "openstack.k-orc.cloud", "controllerKind": "Network"}
2024-11-11T12:09:30Z    INFO    Starting EventSource    {"controller": "image", "controllerGroup": "openstack.k-orc.cloud", "controllerKind": "Image", "source": "kind source: *v1alpha1.Image"}
2024-11-11T12:09:30Z    INFO    Starting Controller     {"controller": "image", "controllerGroup": "openstack.k-orc.cloud", "controllerKind": "Image"}
2024-11-11T12:09:30Z    INFO    Starting workers        {"controller": "image", "controllerGroup": "openstack.k-orc.cloud", "controllerKind": "Image", "worker count": 1}
2024-11-11T12:09:30Z    INFO    Starting workers        {"controller": "network", "controllerGroup": "openstack.k-orc.cloud", "controllerKind": "Network", "worker count": 1}
```

To recompile, kill the process with ++ctrl+c++ and re-run it.

## Create an example ORC resource

### Create OpenStack credentials

Create a `clouds.yaml` file in `examples/credentials`. The name of the cloud in this clouds.yaml must be `openstack`.

This file is in both `.gitignore` and `.dockerignore`, so should not be accidentally added to the git repo or a container build.

We will create an appropriately formatted secret containing these credentials in the next step.

Note that we intentionally create credentials separately from other modules.
This allows us to delete an entire example kustomize module without also
deleting the credentials, which would prevent the deletion from completing.

### Initialise the kustomize environment and load OpenStack credentials

In the examples directory, run:
```bash
$ make
echo "$KUSTOMIZATION" > components/dev-settings/kustomization.yaml
kubectl apply -k apply/credentials --server-side
secret/mbooth-cloud-config-g4ckbm986f serverside-applied
```

This did 2 things. Firstly, it generated the `dev-settings` kustomize component, which adds the current user's username as a `namePrefix`. The purpose of this is to avoid naming conflicts between developers when generating resources in shared clouds, and also to identify culprits if the resources are not cleaned up.

Secondly, it created a secret containing the clouds.yaml we copied into place above. If you missed the first step you will see an error like:
```bash
$ make
Makefile:41: *** You must copy an appropriate clouds.yaml to /home/mbooth/src/openstack-resource-controller/examples/credentials/clouds.yaml. The name of the contained cloud must be 'openstack'..  Stop.
```

Note that failing to initialise the kustomize environment will result in an error like the following when attempting to generate one of the example modules:
```
Error: accumulating components: accumulateDirectory: "couldn't make target for path '.../openstack-resource-controller/examples/components/dev-settings': unable to find one of 'kustomization.yaml', 'kustomization.yml' or 'Kustomization' in directory '.../openstack-resource-controller/examples/components/dev-settings'"
```

### Create an ORC resource

To generate the `managed-network` example:
```bash
$ cd examples/apply/managed-network
$ kubectl apply -k . --server-side
network.openstack.k-orc.cloud/mbooth-orc-managed-network serverside-applied
```

NOTE: This will not create the required secret! To do this, see the section
above on creating credentials.

To cleanup the `managed-network` example:
```bash
$ kubectl delete -k .
```
