# uapi_dhcp_dnsmasq is a singleton: it cannot be created or destroyed. Applying
# writes the settings; destroying only drops it from state.
resource "uapi_dhcp_dnsmasq" "this" {
  domain        = "lan"
  local         = "/lan/"
  authoritative = true
}
