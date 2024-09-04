resource "identitynow_connector_rule" "demo" {
  name        = "test"
  description = "test rule description"
  type        = "BuildMap"
  signature = {
    input = [
      {
        name        = "firstName"
        description = "First Name of identity"
        type        = "String"
      },
      {
        name        = "lastName"
        description = "Last Name of identity"
        type        = "String"
      }
    ]
    output = {
      name      = "fullName"
      description = "Full Name of identity"
      type        = "String"
    }
  }
  source_code = {
    version = "1.0"
    script  = "return firstName + \" \" + lastName;"
  }
  attributes = jsonencode({
    attribute1 = "yes"
    attribute2 = "no"
  })
}
