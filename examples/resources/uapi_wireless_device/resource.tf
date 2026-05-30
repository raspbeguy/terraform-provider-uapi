resource "uapi_wireless_device" "radio0" {
  type    = "mac80211"
  band    = "5g"
  channel = "36"
  htmode  = "VHT80"
  country = "FR"
}
