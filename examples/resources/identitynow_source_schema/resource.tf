resource "identitynow_source_schema" "account_schema" {
  name                = "account"
  source_id           = identitynow_source.demo_source.id
  native_object_type  = "User"
  identity_attribute  = "id"
  display_attribute   = "uid"
  include_permissions = false
  features            = []
  configuration       = jsonencode({})
  attributes = [
    {
      name           = "id"
      type           = "STRING"
      description    = "ID of the user"
      is_multi       = false
      is_entitlement = false
      is_group       = false
    },
    {
      name           = "uid"
      type           = "STRING"
      description    = "UID of the user"
      is_multi       = false
      is_entitlement = false
      is_group       = false
    },
    {
      name           = "firstName"
      type           = "STRING"
      description    = "FIRST NAME"
      is_multi       = false
      is_entitlement = false
      is_group       = false
    },
    {
      name           = "lastName"
      type           = "STRING"
      description    = "LAST NAME"
      is_multi       = false
      is_entitlement = false
      is_group       = false
    }
  ]
}
