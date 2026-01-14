---
page_title: "openai_project_api_key Resource"
description: |-
  Manages an OpenAI project API key.
---

# openai_project_api_key

Manages an OpenAI project API key.

~> **Note:** API keys cannot be created directly via this resource. They are created via service accounts. Use `openai_project_service_account` to create new API keys. This resource is primarily for importing and managing existing keys.

## Example Usage

```hcl
resource "openai_project_api_key" "example" {
  project_id = openai_project.example.id
  id         = "key_abc123"
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the API key.
- `project_id` - (Required) The ID of the project this API key belongs to.

## Attribute Reference

- `name` - The name of the API key.
- `redacted_value` - The redacted value of the API key (e.g., `sk-...abc`).
- `created_at` - The Unix timestamp (in seconds) of when the API key was created.
- `owner_type` - The type of owner of this API key (`user` or `service_account`).
- `owner_id` - The ID of the owner (user or service account).
- `owner_name` - The name of the owner.

## Import

API keys can be imported using the format `project_id/api_key_id`:

```shell
terraform import openai_project_api_key.example proj_abc123/key_xyz789
```
