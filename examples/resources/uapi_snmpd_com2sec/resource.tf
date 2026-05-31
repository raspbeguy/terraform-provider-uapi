resource "uapi_snmpd_com2sec" "ro" {
  secname   = "ro"
  source    = "default"
  community = "public"
}
