# Import by id (the sha256 prefix uapi assigns to the key).
# The key blob is not returned by uapi, so you must set `key` in config after importing.
terraform import uapi_authorized_key.admin 0a1b2c3d4e5f
