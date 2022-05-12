resource "google_storage_bucket" "my_sql_backup_bucket" {
  name          = var.MY_SQL_BACKUP_BUCKET_NAME
  location      = "EU"
  storage_class = "STANDARD"
  force_destroy = true

  lifecycle_rule {
    condition {
      num_newer_versions = 30
    }
    action {
      type = "Delete"
    }
  }
}

resource "kubernetes_namespace" "my_sql_backup_ns" {
  metadata {
    annotations = {
      name = var.NAMESPACE
    }

    labels = {
      mylabel = var.NAMESPACE
    }

    name = var.NAMESPACE
  }
}

resource "google_service_account" "mysql-backup-sa" {
  account_id   = "mysql-backup-sa"
  display_name = "mysql-backup-sa"
}

resource "google_project_iam_member" "mysql-backup-role-membership" {
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.mysql-backup-sa.email}"
}

resource "google_service_account_key" "mysql_backup_key" {
  service_account_id = google_service_account.mysql-backup-sa.name
}

locals {
  jobdeployment = templatefile(
    "${path.module}/templates/my_database_backup.tftpl",
    {
      GCP_BUCKET_NAME = var.MY_SQL_BACKUP_BUCKET_NAME
      MY_SQL_DB_USER  = var.MY_SQL_DB_USER
      MY_SQL_DB_HOST  = var.MY_SQL_DB_HOST
      IMAGE_VERSION   = var.IMAGE_VERSION
    }
  )
}

resource "kubernetes_secret" "google-application-credentials" {
  metadata {
    name = "my-database-backup"
    namespace = var.NAMESPACE
    annotations = {
      "kubernetes.io/service-account.name" = google_service_account.mysql-backup-sa.email
    }
  }
  data = {
    "gcp_gcloud_auth" = base64encode(base64decode(google_service_account_key.mysql_backup_key.private_key))
    "database_password" = var.MY_SQL_DB_PASSWORD
  }

  depends_on = [ kubernetes_namespace.my_sql_backup_ns ]
}

resource "kubectl_manifest" "kubectl_manifest_job" {
    override_namespace = var.NAMESPACE
    provider = kubectl
    yaml_body = local.jobdeployment
    depends_on = [ kubernetes_namespace.my_sql_backup_ns ]
}  