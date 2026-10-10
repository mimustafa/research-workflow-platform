output "cluster_name" {
  description = "Name of the created Kind cluster."
  value       = kind_cluster.research_workflow.name
}

output "kubeconfig_path" {
  description = "Kubeconfig used by kubectl and Helm."
  value       = kind_cluster.research_workflow.kubeconfig_path
}

output "namespace" {
  description = "Namespace reserved for the Helm release."
  value       = kubernetes_namespace_v1.research_workflow.metadata[0].name
}
