resource "uapi_network_device" "br_lan" {
  name  = "br-lan"
  type  = "bridge"
  ports = ["eth0", "eth1"]
}
