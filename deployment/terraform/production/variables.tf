variable "do_token" {
  type = string
}

variable "region" {
  type    = string
  default = "fra1"
}

variable "kubernetes_version" {
  type    = string
  default = "1.22.11-do.0"
}

variable "node_pool_min_nodes" {
  type    = number
  default = 1
}

variable "node_pool_max_nodes" {
  type    = number
  default = 3
}

variable "name" {
  type    = string
  default = "rental"
}

variable "postgres_version" {
  type    = string
  default = "14"
}
