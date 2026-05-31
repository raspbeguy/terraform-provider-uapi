# uapi_prometheus_node_exporter_lua_config is a singleton: it cannot be created
# or destroyed. Applying writes the settings; destroying only drops it from state.
resource "uapi_prometheus_node_exporter_lua_config" "this" {
  listen_port = "9100"
  cpu         = true
  meminfo     = true
  netdev      = true
  loadavg     = true
}
