data "uapi_dhcp_leases6" "current" {}

output "active_leases6" {
  value = data.uapi_dhcp_leases6.current.leases
}
