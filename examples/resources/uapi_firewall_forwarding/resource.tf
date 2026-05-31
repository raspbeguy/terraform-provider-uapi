resource "uapi_firewall_forwarding" "lan_to_wan" {
  src  = "lan"
  dest = "wan"
}
