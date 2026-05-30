resource "uapi_firewall_redirect" "web" {
  name   = "Forward-HTTP"
  target = "DNAT"

  match = {
    src_zone  = "wan"
    proto     = ["tcp"]
    src_dport = ["80"]
    dest_ip   = ["192.168.1.10"]
    dest_port = ["8080"]
  }
}
