resource "google_cloud_run_v2_service" "default" {

  depends_on = [
    google_service_account_iam_member.iam_member,
    google_service_account_iam_binding.act_as_iam,
  ]

  name     = var.service_name
  location = var.cloud_region
  project  = var.project_id

  template {
    containers {
      image = var.image_url

      ports {
        container_port = 7700
      }

      resources {
        limits = {
          memory = "1024Mi"
          cpu    = "2"
        }
      }

      env {
        name  = "MEILI_MASTER_KEY"
        value = var.meilisearch_master_key
      }

      env {
        name  = "MEILI_NO_ANALYTICS"
        value = var.meilisearch_no_analytics
      }

      env {
        name  = "MEILI_ENV"
        value = var.meilisearch_env
      }

      env {
        name  = "TZ"
        value = var.tz
      }
    }

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    service_account = var.service_name
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

}

data "google_project" "project" {
  project_id = var.project_id
}

resource "google_service_account" "default_service_account" {
  account_id   = var.service_name
  display_name = var.service_name
  project      = data.google_project.project.project_id
}

resource "google_service_account_iam_binding" "act_as_iam" {
  service_account_id = google_service_account.default_service_account.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${google_service_account.default_service_account.email}",
  ]
}

resource "google_service_account_iam_member" "iam_member" {
  service_account_id = google_service_account.default_service_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.default_service_account.email}"
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"

    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_v2_service_iam_policy" "noauth" {
  project  = google_cloud_run_v2_service.default.project
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

# gsutil mb -p duckhome-firebase -c STANDARD -l southamerica-east1 gs://duckhome-pps-terraform-state
terraform {
  backend "gcs" {
    bucket = "duckhome-pps-terraform-state"
    prefix = "terraform/state"
  }
}
