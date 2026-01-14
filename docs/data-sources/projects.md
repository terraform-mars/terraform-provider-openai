---
page_title: "openai_projects Data Source"
description: |-
  Fetches a list of all OpenAI projects in the organization.
---

# openai_projects

Fetches a list of all OpenAI projects in the organization.

## Example Usage

```hcl
data "openai_projects" "all" {}

output "project_names" {
  value = [for p in data.openai_projects.all.projects : p.name]
}
```

### Including Archived Projects

```hcl
data "openai_projects" "all" {
  include_archived = true
}
```

## Argument Reference

- `include_archived` - (Optional) Whether to include archived projects. Defaults to `false`.

## Attribute Reference

- `projects` - List of projects in the organization. Each project contains:
  - `id` - The unique identifier of the project.
  - `name` - The name of the project.
  - `status` - The status of the project.
  - `created_at` - The Unix timestamp of when the project was created.
  - `archived_at` - The Unix timestamp of when the project was archived.
