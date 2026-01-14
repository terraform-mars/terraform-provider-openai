---
page_title: "openai_project_service_account Resource"
description: |-
  Creates an OpenAI service account within a project.
---

# openai_project_service_account

Creates an OpenAI service account within a project. Service accounts are used for programmatic API access and are created with an associated API key.

!> **Important:** The `api_key_value` attribute is only available immediately after creation. Store it securely as it cannot be retrieved later.

## Example Usage

```hcl
resource "openai_project" "example" {
  name = "my-project"
}

resource "openai_project_service_account" "example" {
  project_id = openai_project.example.id
  name       = "terraform-service-account"
}

output "api_key" {
  value     = openai_project_service_account.example.api_key_value
  sensitive = true
}
```

## Argument Reference

- `project_id` - (Required) The ID of the project to create the service account in. Forces new resource if changed.
- `name` - (Required) The name of the service account. Forces new resource if changed.

## Attribute Reference

- `id` - The unique identifier of the service account.
- `role` - The role of the service account (currently always `member`).
- `created_at` - The Unix timestamp (in seconds) of when the service account was created.
- `api_key_id` - The ID of the API key associated with this service account.
- `api_key_value` - (Sensitive) The API key value. Only available immediately after creation.

## Import

Service accounts can be imported using the format `project_id/service_account_id`:

```shell
terraform import openai_project_service_account.example proj_abc123/sa_xyz789
```

~> **Note:** When importing, the `api_key_value` will not be available as it is only provided at creation time.
