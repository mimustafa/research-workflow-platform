variable "cluster_name" {
  description = "Name of the local Kind cluster."
  type        = string
  default     = "research-workflow"
}

variable "namespace" {
  description = "Namespace that will contain the Helm release."
  type        = string
  default     = "research-workflow"
}
