variable "ENVIRONMENT_NAME" {
  type  = string
}
variable "K8S_HOST" {
  type  = string
}
variable "K8S_CLIENT_KEY" {
  type  = string
}
variable "K8S_TOKEN" {
  type  = string
}
variable "K8S_CLUSTER_CA_CERTIFICATE" {
  type  = string
}

variable "MY_SQL_BACKUP_BUCKET_NAME" {
  type  = string
}

variable "NAMESPACE" {
  type  = string
}

variable "MY_SQL_DB_USER" {
  type  = string
}

variable "MY_SQL_DB_PASSWORD" {
  type  = string
}

variable "MY_SQL_DB_HOST" {
  type  = string
}

variable "IMAGE_VERSION" {
  type  = string
}