---
page_title: "Provider: OpenAI"
description: |-
  Terraform provider for managing OpenAI organization resources.
---

# OpenAI Provider

The OpenAI provider allows you to manage OpenAI organization resources such as projects, API keys, and service accounts using the [Admin API](https://platform.openai.com/docs/api-reference/administration).

## Example Usage

```hcl
terraform {
  required_providers {
    openai = {
      source  = "terraform-mars/openai"
      version = "~> 0.1"
    }
  }
}

provider "openai" {
  admin_key = var.openai_admin_key
}

resource "openai_project" "example" {
  name = "my-project"
}

resource "openai_project_service_account" "example" {
  project_id = openai_project.example.id
  name       = "terraform-service-account"
}
```

## Authentication

The provider requires an OpenAI Admin API key. You can obtain one from the [OpenAI Platform](https://platform.openai.com/organization/admin-keys).

Configure the key via:
- Provider configuration: `admin_key`
- Environment variable: `OPENAI_ADMIN_KEY`

## Argument Reference

- `admin_key` - (Optional) OpenAI Admin API key. Can also be set via `OPENAI_ADMIN_KEY` environment variable.
- `base_url` - (Optional) OpenAI API base URL. Defaults to `https://api.openai.com/v1`. Can also be set via `OPENAI_BASE_URL` environment variable.
