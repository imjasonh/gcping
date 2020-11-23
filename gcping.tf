////// Providers

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.48.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "3.48.0"
    }
  }
}

provider "google" {
  project = var.project
}

provider "google-beta" {
  project = var.project
}

////// Variables

variable "image" {
  type = string
}

variable "project" {
  type    = string
  default = "gcping-1369"
}

variable "domain" {
  type    = string
  default = "gcping.com"
}

// TODO: generate this
// https://github.com/hashicorp/terraform-provider-google/issues/7850
variable "regions" {
  type        = list(string)
  description = "deploy to regions"

  default = [
    "asia-east1",
    "asia-east2",
    "asia-northeast1",
    "asia-northeast2",
    "asia-northeast3",
    "asia-south1",
    "asia-southeast1",
    "asia-southeast2",
    "australia-southeast1",
    "europe-north1",
    "europe-west1",
    "europe-west2",
    "europe-west3",
    "europe-west4",
    "europe-west6",
    "northamerica-northeast1",
    "southamerica-east1",
    "us-central1",
    "us-east1",
    "us-east4",
    "us-west1"
  ]
}

////// Cloud Run

// Enable Cloud Run API.
resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

// Deploy image to each region.
resource "google_cloud_run_service" "regions" {
  for_each = toset(var.regions)
  name     = each.key
  location = each.key

  metadata {
    annotations = {
      "run.googleapis.com/launch-stage" = "BETA"
    }
  }

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale" = "1"
        "autoscaling.knative.dev/maxScale" = "10"
        "run.googleapis.com/launch-stage"  = "BETA"
      }
    }
    spec {
      containers {
        image = var.image
        env {
          name  = "REGION"
          value = each.key
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.run]
}

// Make each service invokable by all users.
resource "google_cloud_run_service_iam_member" "allUsers" {
  for_each = toset(var.regions)

  service  = google_cloud_run_service.regions[each.key].name
  location = each.key
  role     = "roles/run.invoker"
  member   = "allUsers"

  depends_on = [google_cloud_run_service.regions]
}

// Create a regional network endpoint group (NEG) for each regional Cloud Run service.
resource "google_compute_region_network_endpoint_group" "regions" {
  for_each = toset(var.regions)

  name                  = each.key
  network_endpoint_type = "SERVERLESS"
  region                = each.key
  cloud_run {
    service = google_cloud_run_service.regions[each.key].name
  }
}

////// Regional Domain + Load Balancer config

resource "google_compute_global_forwarding_rule" "regions" {
  for_each = toset(var.regions)

  name       = each.key
  target     = google_compute_target_https_proxy.regions[each.key].id
  port_range = "443"
}

// Print regional LB IP addresses.
output "regions" {
  value = {
    for fwd in google_compute_global_forwarding_rule.regions :
    fwd.name => fwd.ip_address
  }
}

resource "google_compute_managed_ssl_certificate" "regions" {
  for_each = toset(var.regions)
  provider = google-beta

  name = each.key
  managed {
    domains = ["${each.key}.${var.domain}."]
  }
}

resource "google_compute_target_https_proxy" "regions" {
  for_each = toset(var.regions)
  provider = google-beta

  name             = each.key
  url_map          = google_compute_url_map.regions[each.key].id
  ssl_certificates = [google_compute_managed_ssl_certificate.regions[each.key].id]
}

resource "google_compute_url_map" "regions" {
  for_each = toset(var.regions)
  provider = google-beta

  name            = each.key
  description     = "a description"
  default_service = google_compute_backend_service.regions[each.key].id
}

resource "google_compute_backend_service" "regions" {
  for_each = toset(var.regions)

  name       = each.key
  enable_cdn = true
  backend {
    group = google_compute_region_network_endpoint_group.regions[each.key].id
  }
}

////// Global Domain + Load Balancer config

resource "google_compute_global_forwarding_rule" "global" {
  name       = "global"
  target     = google_compute_target_https_proxy.global.id
  port_range = "443"
}

// Print global LB IP address.
output "global" {
  value = google_compute_global_forwarding_rule.global.ip_address
}

resource "google_compute_managed_ssl_certificate" "global" {
  provider = google-beta

  name = "global"
  managed {
    domains = [
      "global.${var.domain}",
      "www.${var.domain}",
      var.domain,
    ]
  }
}

resource "google_compute_target_https_proxy" "global" {
  provider = google-beta

  name             = "global"
  url_map          = google_compute_url_map.global.id
  ssl_certificates = [google_compute_managed_ssl_certificate.global.id]
}

resource "google_compute_url_map" "global" {
  provider = google-beta

  name            = "global"
  description     = "a description"
  default_service = google_compute_backend_service.global.id
}

// Create a global backend service with a backend for each regional NEG.
resource "google_compute_backend_service" "global" {
  name       = "global"
  enable_cdn = true

  // Add a backend for each regional NEG.
  dynamic "backend" {
    for_each = google_compute_region_network_endpoint_group.regions
    content {
      group = backend.value["id"]
    }
  }
}
