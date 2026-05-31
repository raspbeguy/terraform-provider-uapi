resource "uapi_network_interface" "lan" {
  device  = "br-lan"
  proto   = "static"
  ipaddr  = "192.168.1.1"
  netmask = "255.255.255.0"
  dns     = ["1.1.1.1", "9.9.9.9"]
}

# A WireGuard interface. private_key is write-only (never read back).
resource "uapi_network_interface" "wg0" {
  proto       = "wireguard"
  private_key = var.wg_private_key
  listen_port = "51820"
  addresses   = ["10.10.0.1/24"]
}

variable "wg_private_key" {
  type      = string
  sensitive = true
}
