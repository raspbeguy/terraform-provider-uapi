resource "uapi_network_route" "to_dmz" {
  interface = "lan"
  target    = "10.10.0.0/24"
  gateway   = "192.168.1.254"
  metric    = "10"
}
