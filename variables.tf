variable "project_id" {
  description = "The ID of the project to apply any resources to."
  type        = string
}

variable "cloud_region" {
  description = "The region to deploy to."
  type        = string
}

variable "zone" {
  description = "Cloud zone"
  type        = string
}

variable "service_name" {
  description = "Name of the Cloud Run service."
  type        = string
}

variable "meilisearch_image_url" {
  description = "URL of the MeiliSearch Docker image to deploy."
  type        = string
}

variable "meilisync_image_url" {
  description = "URL of the MeiliSync Docker image to deploy."
  type        = string
}

variable "meilisearch_host" {
  description = "Server Hostname."
  type        = string
}

variable "meilisearch_master_key" {
  description = "Master key."
  type        = string
}

variable "meilisearch_no_analytics" {
  description = "Disable analytics."
  type        = string
}

variable "meilisearch_env" {
  description = "Environment."
  type        = string
}

variable "meilisearch_bucket_name" {
  description = "MeiliSearch bucket name."
  type        = string
}

variable "meilisync_bucket_name" {
  description = "MeiliSync bucket name."
  type        = string
}

variable "tz" {
  description = "Timezone."
  type        = string
}
