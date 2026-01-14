---
page_title: "openai_project Resource"
description: |-
  Manages an OpenAI project.
---

# openai_project

Manages an OpenAI project. Projects are used to organize API keys and resources within an organization.

~> **Note:** OpenAI projects cannot be deleted, only archived. When this resource is destroyed, the project will be archived.

## Example Usage

```hcl
resource "openai_project" "example" {
  name = "my-project"
}
```

## Argument Reference

- `name` - (Required) The name of the project. This appears in reporting.

## Attribute Reference

- `id` - The unique identifier of the project.
- `status` - The status of the project (`active` or `archived`).
- `created_at` - The Unix timestamp (in seconds) of when the project was created.
- `archived_at` - The Unix timestamp (in seconds) of when the project was archived, or null if not archived.

## Import

Projects can be imported using the project ID:

```shell
terraform import openai_project.example proj_abc123
```
