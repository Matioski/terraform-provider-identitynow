resource "identitynow_transform" "demo_transform" {
  name = "demo transform calculate"
  type = "substring"
  attributes = jsonencode({
    input = {
      attributes = {
        attributeName = "firstName"
        sourceName    = identitynow_source.demo_source.name
      }
      type = "accountAttribute"
    }
    end = -1.0
    begin = {
      attributes = {
        substring = ","
      }
      type = "indexOf"
    }
    beginOffset = 3.0
  })
}
