terraform {
  backend "s3" {
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    endpoint                    = "https://fra1.digitaloceanspaces.com"
    // We need a valid AWS region even though we're using DigitalOcean Spaces.
    region = "us-east-1"
    bucket = "tf-rental"
    key    = "terraform.tfstate"
  }
}
