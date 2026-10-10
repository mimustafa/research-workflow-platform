terraform {
  required_version = ">= 1.8.0"

  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = "~> 0.8"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.36"
    }
  }
}

provider "kind" {}

resource "kind_cluster" "research_workflow" {
  name            = var.cluster_name
  wait_for_ready  = true
  kubeconfig_path = pathexpand("~/.kube/${var.cluster_name}.yaml")

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role = "control-plane"
    }
  }
}

provider "kubernetes" {
  config_path = kind_cluster.research_workflow.kubeconfig_path
}

resource "kubernetes_namespace_v1" "research_workflow" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "research-workflow"
    }
  }
}
