resource "uapi_wireless_interface" "home" {
  device     = uapi_wireless_device.radio0.id
  network    = "lan"
  mode       = "ap"
  ssid       = "home-net"
  encryption = "psk2"
  key        = var.wifi_key # write-only: never returned by the API
}

variable "wifi_key" {
  type      = string
  sensitive = true
}
