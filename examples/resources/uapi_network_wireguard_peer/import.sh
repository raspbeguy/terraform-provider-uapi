# A uapi-managed WireGuard peer is imported by its stable id.
terraform import uapi_network_wireguard_peer.laptop g_01HX0000000000000000000000

# Importing a pre-existing anonymous (unmanaged) section adopts it: uapi renames
# it to a stable id and the provider emits a warning naming the old and new ids.
terraform import uapi_network_wireguard_peer.laptop cfg0a1b2c
