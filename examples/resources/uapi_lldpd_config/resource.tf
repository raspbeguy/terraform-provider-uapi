# uapi_lldpd_config is a singleton: it cannot be created or destroyed. Applying
# writes the settings; destroying only drops it from state.
resource "uapi_lldpd_config" "this" {
  enable_lldpmed    = true
  lldp_description   = true
  lldp_capabilities = true
  interface         = ["lan"]
}
