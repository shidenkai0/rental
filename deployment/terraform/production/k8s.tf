resource "digitalocean_kubernetes_cluster" "default" {
  name     = format("%s-prod", var.name)
  region   = var.region
  version  = var.kubernetes_version
  vpc_uuid = digitalocean_vpc.default.id
  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-2gb"
    auto_scale = true
    min_nodes  = var.node_pool_min_nodes
    max_nodes  = var.node_pool_max_nodes
  }

}


output "kube_config_base64" {
  sensitive = true
  value     = base64encode(digitalocean_kubernetes_cluster.default.kube_config.0.raw_config)
}
