resource "uapi_network_rule" "from_guests" {
  src      = "192.168.3.0/24"
  priority = "100"
  lookup   = "200"
}
