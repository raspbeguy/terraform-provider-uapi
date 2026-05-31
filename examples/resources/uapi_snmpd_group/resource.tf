resource "uapi_snmpd_group" "public" {
  group   = "public"
  version = "v2c"
  secname = uapi_snmpd_com2sec.ro.secname
}
