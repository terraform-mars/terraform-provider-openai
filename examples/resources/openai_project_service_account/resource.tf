# Create a service account for CI/CD
resource "openai_project_service_account" "github_actions" {
  project_id = openai_project.production.id
  name       = "github-actions"
}

# The API key is available immediately after creation
output "github_actions_api_key" {
  value     = openai_project_service_account.github_actions.api_key_value
  sensitive = true
}

# Create service accounts for each environment
resource "openai_project_service_account" "backend" {
  for_each = {
    dev     = openai_project.development.id
    prod    = openai_project.production.id
  }

  project_id = each.value
  name       = "backend-${each.key}"
}
