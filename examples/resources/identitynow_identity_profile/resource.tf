resource "identitynow_identity_profile" "demo_identity_profile" {
  name        = "DEMO Terraform Identity Profile"
  description = "DEMO Terraform Identity Profile"
  authoritative_source = {
    id = identitynow_source.demo_source.id
  }
  owner = {
    id   = data.identitynow_identity.default_owner.id
    name = data.identitynow_identity.default_owner.name
  }
  identity_attribute_config = {
    enabled = false
    attribute_transforms = [
      {
        identity_attribute_name = "uid"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = identitynow_source.demo_source.name
            attributeName = "uid"
            sourceId      = identitynow_source.demo_source.id
          })
        }
      },
      {
        identity_attribute_name = "email"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = identitynow_source.demo_source.name
            attributeName = "uid"
            sourceId      = identitynow_source.demo_source.id
          })
        }
      },

      {
        identity_attribute_name = "lastname"
        transform_definition = {
          type = "accountAttribute"
          attributes = jsonencode({
            sourceName    = identitynow_source.demo_source.name
            attributeName = "uid"
            sourceId      = identitynow_source.demo_source.id
          })
        }
      },
      {
        identity_attribute_name = "firstname"
        transform_definition = {
          type = "reference"
          attributes = jsonencode({
            input = {
              attributes = {
                attributeName = "firstName"
                sourceName    = identitynow_source.demo_source.name
                sourceId      = identitynow_source.demo_source.id
              }
              type = "accountAttribute"
            }
            id = identitynow_transform.demo_transform.name
          })
        }
      }
    ]
  }
}
