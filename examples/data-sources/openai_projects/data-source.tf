# List all active projects
data "openai_projects" "all" {}

output "project_names" {
  value = [for p in data.openai_projects.all.projects : p.name]
}

# Include archived projects
data "openai_projects" "all_including_archived" {
  include_archived = true
}
