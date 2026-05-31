# uapi_firewall_defaults is a singleton: it cannot be created or destroyed.
# Applying writes the settings; destroying only drops it from state.
resource "uapi_firewall_defaults" "this" {
  input   = "ACCEPT"
  output  = "ACCEPT"
  forward = "REJECT"
}
