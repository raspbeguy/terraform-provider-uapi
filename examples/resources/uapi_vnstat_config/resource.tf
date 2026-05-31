# uapi_vnstat_config is a singleton: it cannot be created or destroyed. Applying
# writes the settings; destroying only drops it from state.
resource "uapi_vnstat_config" "this" {
  database_dir         = "/var/lib/vnstat"
  interface_5min_hours = "48"
  month_rotate         = "1"
}
