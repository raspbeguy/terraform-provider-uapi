data "uapi_authorized_key" "admin" {
  id = "0a1b2c3d4e5f"
}

output "admin_key_type" {
  value = data.uapi_authorized_key.admin.type
}
