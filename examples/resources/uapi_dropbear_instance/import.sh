# A uapi-managed dropbear instance is imported by its stable id.
terraform import uapi_dropbear_instance.lan d_01HX0000000000000000000000

# Importing a pre-existing anonymous (unmanaged) section adopts it: uapi renames
# it to a stable id and the provider emits a warning naming the old and new ids.
terraform import uapi_dropbear_instance.lan cfg0a1b2c
