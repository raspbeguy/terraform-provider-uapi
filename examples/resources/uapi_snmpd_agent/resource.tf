resource "uapi_snmpd_agent" "this" {
  agentaddress = ["UDP:161", "udp6:161"]
}
