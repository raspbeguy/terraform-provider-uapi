resource "uapi_sqm_queue" "wan" {
  interface = "wan"
  download  = "100000"
  upload    = "20000"
  qdisc     = "cake"
  script    = "piece_of_cake.qos"
  linklayer = "ethernet"
}
