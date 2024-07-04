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

variable "image_url" {
  description = "URL of the Docker image to deploy."
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

variable "tz" {
  description = "Timezone."
  type        = string
}
