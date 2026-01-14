# Terraform Provider for OpenAI

A Terraform provider for managing OpenAI organization resources like projects, service accounts, and API keys.

## Features

This provider focuses on **infrastructure management** - the resources enterprises need to manage at scale:

- **Projects** - Create and manage organizational projects
- **Service Accounts** - Programmatic access with auto-generated API keys
- **API Keys** - Manage existing API keys

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for building)
- OpenAI Admin API Key ([create one here](https://platform.openai.com/settings/organization/admin-keys))

## Installation

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
  # Can also use OPENAI_ADMIN_KEY environment variable
  admin_key = var.openai_admin_key
}
```

## Usage

### Create Projects

```hcl
resource "openai_project" "production" {
  name = "production"
}

resource "openai_project" "development" {
  name = "development"
}
```

### Create Service Accounts with API Keys

```hcl
resource "openai_project_service_account" "backend" {
  project_id = openai_project.production.id
  name       = "backend-service"
}

# The API key value is only available at creation time
output "api_key" {
  value     = openai_project_service_account.backend.api_key_value
  sensitive = true
}
```

### Multi-Environment Setup

```hcl
# Create projects for each environment
resource "openai_project" "envs" {
  for_each = toset(["dev", "staging", "prod"])
  name     = each.key
}

# Create service accounts for each project
resource "openai_project_service_account" "backend" {
  for_each   = openai_project.envs
  project_id = each.value.id
  name       = "backend-${each.key}"
}

# Export keys for CI/CD
output "api_keys" {
  sensitive = true
  value = {
    for env, sa in openai_project_service_account.backend :
    env => sa.api_key_value
  }
}
```

### List Existing Projects

```hcl
data "openai_projects" "all" {}

output "projects" {
  value = data.openai_projects.all.projects
}
```

## Resources

| Resource | Description |
|----------|-------------|
| `openai_project` | Manages an OpenAI project |
| `openai_project_service_account` | Creates a service account with API key |
| `openai_project_api_key` | Manages an existing API key |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| `openai_project` | Fetches a project by ID |
| `openai_projects` | Lists all projects |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `OPENAI_ADMIN_KEY` | OpenAI Admin API key |
| `OPENAI_BASE_URL` | API base URL (default: `https://api.openai.com/v1`) |

## Building

```bash
go build -o terraform-provider-openai
```

## Testing

```bash
make test        # Unit tests
make testacc     # Acceptance tests (requires OPENAI_ADMIN_KEY)
```

## License

MPL-2.0
