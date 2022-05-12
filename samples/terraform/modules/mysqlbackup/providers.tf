
terraform {
  
  required_providers {
    kubectl = {
      source = "gavinbunney/kubectl"
      version = "1.10.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 1.0.0"
    }
  }
}

provider "kubectl" {
  host                   = var.K8S_HOST
  client_key             = var.K8S_CLIENT_KEY
  cluster_ca_certificate = var.K8S_CLUSTER_CA_CERTIFICATE
  token                  = var.K8S_TOKEN
}

provider "kubernetes" {
  host                   = var.K8S_HOST
  client_key             = var.K8S_CLIENT_KEY
  cluster_ca_certificate = var.K8S_CLUSTER_CA_CERTIFICATE
  token                  = var.K8S_TOKEN
}