resource "uapi_network_bridge_vlan" "guests" {
  device = "br-lan"
  vlan   = "30"
  ports  = ["lan1:t", "lan2:u"]
}
