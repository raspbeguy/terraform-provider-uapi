resource "uapi_dropbear_instance" "lan" {
  port          = "22"
  interface     = "lan"
  password_auth = false
  root_login    = false
}
