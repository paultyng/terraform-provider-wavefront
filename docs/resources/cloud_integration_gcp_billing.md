---
layout: "wavefront"
page_title: "Wavefront: "
description: |-
  Provides a Wavefront Cloud Integration for Google Cloud Billing. This allows GCP Billing cloud integrations to be created,
  updated, and deleted.
---

# Resource: wavefront_cloud_integration_gcp_billing

Provides a Wavefront Cloud Integration for Google Cloud Billing. This allows GCP Billing cloud integrations to be created,
updated, and deleted.

## Example usage

```hcl
resource "wavefront_cloud_integration_gcp_billing" "gcp_billing" {
  name       = "Test Integration"
  project_id = "example-gcp-project"
  api_key    = "example-api-key"
  json_key   = <<EOF
{...your gcp key ...}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) A value denoting which cloud service this service integrates with.
* `name` - (Required) The human-readable name of this integration.
* `additional_tags` - (Optional) A list of point tag key-values to add to every point ingested using this integration.
* `force_save` - (Optional) Forces this resource to save, even if errors are present.
* `service_refresh_rate_in_minutes` - (Optional) How often, in minutes, to refresh the service.
* `project_id` - (Required) The Google Cloud Platform (GCP) Project ID.
* `api_key` - (Required) API key for Google Cloud Platform (GCP).
* `json_key` - (Required) Private key for a Google Cloud Platform (GCP) service account within your project.
  The account must have at least Viewer permissions. This key must be in the JSON format generated by GCP.

### Example

```hcl
resource "wavefront_cloud_integration_gcp_billing" "gcp_billing" {
  name            = "Test Integration"
  force_save      = true
  additional_tags = {
    "tag1" = "value1"
    "tag2" = "value2"
  }
  project_id                      = "example-gcp-project"
  api_key                         = "example-api-key"
  json_key                        = <<EOF
{...your gcp key ...}
EOF
  service_refresh_rate_in_minutes = 10
}
```

## Import

GCP Billing Cloud Integrations can be imported by using the `id`, e.g.:

```
$ terraform import wavefront_cloud_integration_gcp_billing.gcp_billing a411c16b-3cf7-4f03-bf11-8ca05aab898d
```