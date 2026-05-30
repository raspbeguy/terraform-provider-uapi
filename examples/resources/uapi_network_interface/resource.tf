resource "uapi_network_interface" "lan" {
  device  = "br-lan"
  proto   = "static"
  ipaddr  = "192.168.1.1"
  netmask = "255.255.255.0"
  dns     = ["1.1.1.1", "9.9.9.9"]
}
