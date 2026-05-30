# uapi_system is a singleton: it cannot be created or destroyed. Applying writes
# the settings; destroying only drops it from state.
resource "uapi_system" "this" {
  hostname = "edge-router"
  timezone = "CET-1CEST,M3.5.0,M10.5.0/3"
  zonename = "Europe/Paris"
}
