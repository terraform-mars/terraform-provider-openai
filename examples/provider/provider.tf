terraform {
  required_providers {
    openai = {
      source  = "terraform-mars/openai"
      version = "~> 0.1"
    }
  }
}

provider "openai" {
  # Admin API key - can also use OPENAI_ADMIN_KEY environment variable
  # admin_key = "sk-admin-..."
}
