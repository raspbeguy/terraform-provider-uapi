# uapi_snmpd_system is a singleton: it cannot be created or destroyed. Applying
# writes the settings; destroying only drops it from state.
resource "uapi_snmpd_system" "this" {
  sys_location = "Server Room A"
  sys_contact  = "noc@example.com"
  sys_name     = "edge-router"
}
