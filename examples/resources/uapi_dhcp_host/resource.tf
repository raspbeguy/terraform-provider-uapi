resource "uapi_dhcp_host" "printer" {
  name = "printer"
  mac  = "aa:bb:cc:dd:ee:ff"
  ip   = "192.168.1.50"
  dns  = true
}
