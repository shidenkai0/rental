resource "digitalocean_database_cluster" "default" {
  name                 = format("%s-prod", var.name)
  engine               = "pg"
  version              = var.postgres_version
  size                 = "db-s-1vcpu-1gb"
  region               = var.region
  node_count           = 1
  private_network_uuid = digitalocean_vpc.default.id
}

resource "digitalocean_database_db" "app" {
  cluster_id = digitalocean_database_cluster.default.id
  name       = var.name
}

output "database_uri" {
  sensitive = true
  value     = digitalocean_database_cluster.default.private_uri
}

resource "digitalocean_database_user" "app" {
  cluster_id = digitalocean_database_cluster.default.id
  name       = "rental"
}

output "database_user_password" {
  sensitive = true
  value     = digitalocean_database_user.app.password
}
