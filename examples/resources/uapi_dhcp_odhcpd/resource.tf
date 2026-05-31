# uapi_dhcp_odhcpd is a singleton: it cannot be created or destroyed. Applying
# writes the settings; destroying only drops it from state.
resource "uapi_dhcp_odhcpd" "this" {
  maindhcp = false
  loglevel = "4"
}
