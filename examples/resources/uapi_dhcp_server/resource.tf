resource "uapi_dhcp_server" "lan" {
  interface = "lan"
  start     = "100"
  limit     = "150"
  leasetime = "12h"
}
