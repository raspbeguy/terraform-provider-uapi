data "uapi_dhcp_leases" "current" {}

output "active_leases" {
  value = data.uapi_dhcp_leases.current.leases
}
