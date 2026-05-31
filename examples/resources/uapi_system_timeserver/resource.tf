resource "uapi_system_timeserver" "ntp" {
  enabled  = true
  use_dhcp = false
  server = [
    "0.openwrt.pool.ntp.org",
    "1.openwrt.pool.ntp.org",
  ]
}
