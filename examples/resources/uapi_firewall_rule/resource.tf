resource "uapi_firewall_rule" "allow_ssh_from_wan" {
  name    = "Allow-SSH-from-WAN"
  target  = "ACCEPT"
  enabled = true

  match = {
    src_zone  = "wan"
    proto     = ["tcp"]
    dest_port = ["22"]
  }
}
