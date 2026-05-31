# A uapi-managed snmpd com2sec is imported by its stable id.
terraform import uapi_snmpd_com2sec.ro c_01HX0000000000000000000000

# Importing a pre-existing anonymous (unmanaged) section adopts it: uapi renames
# it to a stable id and the provider emits a warning naming the old and new ids.
terraform import uapi_snmpd_com2sec.ro cfg0a1b2c
