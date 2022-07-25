resource "digitalocean_container_registry" "rental" {
  name                   = "mmess"
  region                 = var.region
  subscription_tier_slug = "basic"
}

output "registry_url" {
  value = digitalocean_container_registry.rental.endpoint
}
