# Look up an existing project by ID
data "openai_project" "existing" {
  id = "proj_abc123"
}

output "project_name" {
  value = data.openai_project.existing.name
}
