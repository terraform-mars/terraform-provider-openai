# Create a project for production workloads
resource "openai_project" "production" {
  name = "production"
}

# Create a project for development
resource "openai_project" "development" {
  name = "development"
}

output "production_project_id" {
  value = openai_project.production.id
}
