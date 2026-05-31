# A uapi-managed forwarding is imported by its stable id.
terraform import uapi_firewall_forwarding.lan_to_wan f_01HX0000000000000000000000

# Importing a pre-existing anonymous (unmanaged) section adopts it: uapi renames
# it to a stable id and the provider emits a warning naming the old and new ids.
terraform import uapi_firewall_forwarding.lan_to_wan cfg0a1b2c
