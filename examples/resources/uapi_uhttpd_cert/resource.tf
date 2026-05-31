resource "uapi_uhttpd_cert" "router" {
  commonname   = "router.lan"
  days         = "3650"
  bits         = "2048"
  organization = "Example Org"
  country      = "FR"
}
