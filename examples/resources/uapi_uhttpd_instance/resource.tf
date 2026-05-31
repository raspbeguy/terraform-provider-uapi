resource "uapi_uhttpd_instance" "main" {
  listen_http  = ["0.0.0.0:80", "[::]:80"]
  listen_https = ["0.0.0.0:443", "[::]:443"]
  home         = "/www"
  cert         = "/etc/uhttpd.crt"
  key          = "/etc/uhttpd.key"
  index_page   = ["index.html"]
  no_dirlists  = true
}
