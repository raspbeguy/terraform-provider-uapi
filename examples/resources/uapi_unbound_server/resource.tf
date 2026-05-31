# uapi_unbound_server is a singleton: it cannot be created or destroyed. Applying
# writes the settings; destroying only drops it from state.
resource "uapi_unbound_server" "this" {
  enabled        = true
  listen_port    = "53"
  dhcp_link      = "dnsmasq"
  dnssec_enabled = true
  recursion      = "default"
  resource       = "medium"
  protocol       = "mixed"
}
