# A uapi-managed SSID is imported by its stable id. The key is write-only and is
# never read back, so set it in config after import to be able to rotate it.
terraform import uapi_wireless_interface.home f_01HX0000000000000000000000

# Importing a pre-existing anonymous (unmanaged) section adopts it: uapi renames
# it to a stable id and the provider emits a warning naming the old and new ids.
terraform import uapi_wireless_interface.home cfg0a1b2c
