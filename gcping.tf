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
  default = "global.gcping.com."
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

// Deploy ${image} to each region.
resource "google_cloud_run_service" "default" {
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

  service  = google_cloud_run_service.default[each.key].name
  location = each.key
  role     = "roles/run.invoker"
  member   = "allUsers"

  depends_on = [google_cloud_run_service.default]
}

// Print regional Cloud Run service URLs.
output "urls" {
  value = {
    for svc in google_cloud_run_service.default :
    svc.name => svc.status[0].url
  }
}

////// Load Balancer config

// Print the global LB's IP address, to plug into DNS.
output "global-ip" {
  value = google_compute_global_forwarding_rule.default.ip_address
}

resource "google_compute_global_forwarding_rule" "default" {
  name       = "global-rule"
  target     = google_compute_target_https_proxy.default.id
  port_range = "443"
}

resource "google_compute_managed_ssl_certificate" "default" {
  provider = google-beta

  name = "ssl-cert"

  managed {
    domains = [var.domain]
  }
}

resource "google_dns_managed_zone" "zone" {
  provider = google-beta

  name     = "dnszone"
  dns_name = var.domain
}

resource "google_dns_record_set" "set" {
  provider = google-beta

  name         = var.domain
  type         = "A"
  ttl          = 3600
  managed_zone = google_dns_managed_zone.zone.name
  rrdatas      = [google_compute_global_forwarding_rule.default.ip_address]
}

resource "google_compute_target_https_proxy" "default" {
  provider = google-beta

  name             = "https-proxy"
  url_map          = google_compute_url_map.default.id
  ssl_certificates = [google_compute_managed_ssl_certificate.default.id]
}

resource "google_compute_url_map" "default" {
  provider = google-beta

  name            = "url-map"
  description     = "a description"
  default_service = google_compute_backend_service.default.id
}

// Create a regional network endpoint group (NEG) for each regional Cloud Run service.
resource "google_compute_region_network_endpoint_group" "default" {
  for_each = toset(var.regions)

  name                  = "default"
  network_endpoint_type = "SERVERLESS"
  region                = each.key
  cloud_run {
    service = google_cloud_run_service.default[each.key].name
  }
}

// Create a global backend service with a backend for each regional NEG.
resource "google_compute_backend_service" "default" {
  name       = "cloudrun-backend-service" // TODO: rename
  enable_cdn = true

  // Add a backend for each regional NEG.
  dynamic "backend" {
    for_each = google_compute_region_network_endpoint_group.default
    content {
      group = backend.value["id"]
    }
  }
}
