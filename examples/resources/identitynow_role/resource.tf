resource "identitynow_role" "test" {
  name        = "Test Role"
  description = "Creating from terraform"
  owner = {
    id   = data.identitynow_identity.default_owner.id
    name = data.identitynow_identity.default_owner.name
  }
  entitlements = [
  ]
  access_profiles = []
  membership = {
    type = "STANDARD"
    criteria = {
      operation   = "AND"
      key         = null
      string_value = null
      children = [
        {
          operation   = "OR"
          key         = null
          string_value = null
          children = [
            {
              operation = "EQUALS"
              key = {
                type     = "IDENTITY"
                property = "attribute.cloudLifecycleState"
                sourceId = null
              }
              string_value = "Initiated"
              children    = null
            },
            {
              operation = "EQUALS"
              key = {
                type     = "IDENTITY"
                property = "attribute.cloudLifecycleState"
                sourceId = null
              }
              string_value = "Active"
              children    = null
            }
          ]
        },
        {
          operation   = "OR"
          key         = null
          string_value = null
          children = [
            {
              operation = "EQUALS"
              key = {
                type     = "IDENTITY"
                property = "attribute.type"
                sourceId = null
              }
              string_value = "INTERNAL"
              children    = null
            },
            {
              operation = "EQUALS"
              key = {
                type     = "IDENTITY"
                property = "attribute.type"
                sourceId = null
              }
              string_value = "EXTERNAL"
              children    = null
            }
          ]
        }
      ]
    },
    identities = null
  }
}