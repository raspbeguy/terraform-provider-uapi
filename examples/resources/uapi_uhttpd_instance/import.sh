# A uapi-managed uhttpd instance is imported by its stable id.
terraform import uapi_uhttpd_instance.main main

# Importing a pre-existing anonymous (unmanaged) section adopts it: uapi renames
# it to a stable id and the provider emits a warning naming the old and new ids.
terraform import uapi_uhttpd_instance.main cfg0a1b2c
