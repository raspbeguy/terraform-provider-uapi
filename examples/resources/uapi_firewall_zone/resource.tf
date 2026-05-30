resource "uapi_firewall_zone" "dmz" {
  name    = "dmz"
  input   = "DROP"
  output  = "ACCEPT"
  forward = "DROP"
  network = ["dmz"]
  masq    = true
}
