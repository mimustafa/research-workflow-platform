# Local Kubernetes infrastructure

This Terraform example provisions a **local Kind cluster** and a Kubernetes namespace.
It does not create cloud resources, so it is safe to run without a cloud account or
ongoing infrastructure cost.

## Prerequisites

- Docker
- Kind
- Terraform

## Use

```bash
cd infra/terraform
terraform init
terraform apply
export KUBECONFIG="$(terraform output -raw kubeconfig_path)"
kubectl get nodes
```

Build and load the application images before installing the Helm chart:

```bash
docker compose build
kind load docker-image research-workflow-platform-api:latest \
  research-workflow-platform-worker:latest \
  research-workflow-platform-web:latest --name research-workflow
```

Destroy the local cluster and namespace when finished:

```bash
terraform destroy
```
