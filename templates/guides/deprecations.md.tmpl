---
page_title: "Deprecations"
subcategory: "Guides"
description: |-
  Provider attributes that are deprecated but still accepted during a deprecation window.
---

# Deprecations

Attributes here still work but are scheduled for removal in a future major. The
provider mirrors uapi's deprecation policy: a deprecation lands in a minor (both
old and new forms accepted), removal happens no sooner than the next major.

## Active

| Attribute | Replaced by | Deprecated since | Removal target | Migration |
|---|---|---|---|---|
| `uapi_network_interface.name` | `uapi_network_interface.id` | provider 2.2.0 (uapi 2.2.0) | v3 | Set `id` instead of `name` at create. Both pick the uci section name and are accepted during the window; if you supply both they must match. `id` is the universal section-name input on every collection resource (see the "Managing pre-existing named sections" guide); `name` was a 2.1.0-era shim that only worked on `uapi_network_interface`. |

## Migrating `name` to `id`

You do not have to migrate during v2: `name` keeps working for the whole window.
When you do switch to `id`, set it to the value `name` had:

```hcl
resource "uapi_network_interface" "wg0" {
  # name = "wg0"   # was
  id          = "wg0" # now
  proto       = "wireguard"
  private_key = var.wg_key
}
```

Note this **forces one replacement**: both `name` and `id` are create-only, so
dropping `name` is a create-only change and Terraform plans a destroy + create
(`terraform plan` shows the resource must be replaced). The new section keeps the
same name (`id` equals the old `name`), so references to it by `id` are unaffected
afterwards, but the section itself is recreated. For an interface that backs your
management connection, do this in a maintenance window (or leave `name` in place
until then, since it remains valid through v2).
