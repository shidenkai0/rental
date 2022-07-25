resource "digitalocean_vpc" "default" {
  name     = format("%s-prod", var.name)
  region   = var.region
  ip_range = "10.11.10.0/24"
}
