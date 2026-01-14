---
page_title: "openai_project Data Source"
description: |-
  Fetches information about an existing OpenAI project.
---

# openai_project

Fetches information about an existing OpenAI project.

## Example Usage

```hcl
data "openai_project" "example" {
  id = "proj_abc123"
}

output "project_name" {
  value = data.openai_project.example.name
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the project.

## Attribute Reference

- `name` - The name of the project.
- `status` - The status of the project (`active` or `archived`).
- `created_at` - The Unix timestamp (in seconds) of when the project was created.
- `archived_at` - The Unix timestamp (in seconds) of when the project was archived, or null if not archived.
