resource "identitynow_identity_attribute" "demo" {
  name         = "test"
  display_name = "test"
  type         = "string"
  sources = [
    {
      type = "rule"
      properties = jsonencode({
        ruleType = "IdentityAttribute"
        ruleName = "Cloud Promote Identity Attribute"
      })
    }
  ]
}
