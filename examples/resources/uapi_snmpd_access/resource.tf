resource "uapi_snmpd_access" "public" {
  group   = uapi_snmpd_group.public.group
  version = "v2c"
  level   = "noauth"
  prefix  = "exact"
  read    = "all"
}
