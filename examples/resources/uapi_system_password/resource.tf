resource "uapi_system_password" "root" {
  user = "root"

  # Write-only: never stored in state. Sourced here from a variable.
  password_wo = var.root_password

  # Bump this whenever you want the password re-applied.
  password_wo_version = "1"
}

variable "root_password" {
  type      = string
  sensitive = true
}
