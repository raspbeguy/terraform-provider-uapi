resource "uapi_network_wireguard_peer" "laptop" {
  interface            = "wg0"
  description          = "Field laptop"
  public_key           = "xTIBA5rboUvnH4htodjb6e697QjLERt1NAB4mZqp8Dg="
  allowed_ips          = ["10.0.0.2/32"]
  endpoint_host        = "vpn.example.com"
  endpoint_port        = "51820"
  persistent_keepalive = "25"
}
