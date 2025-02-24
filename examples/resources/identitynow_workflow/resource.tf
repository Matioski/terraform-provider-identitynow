resource "identitynow_workflow" "test" {
  name        = "test workflow"
  description = "test workflow from terraform"
  owner = {
    id = "ownerId"
  }
  enabled = false
  definition = {
    start = "Get List of Identities"
    steps = jsonencode({
      "Get List of Identities" = {
        displayName   = "Get List of Identities"
        type          = "action"
        actionId      = "sp:get-identities"
        versionNumber = 2
        attributes = {
          inputQuery = "*"
          searchBy   = "searchQuery"
        }
        nextStep = "End Step"
      }
      "End Step" = {
        displayName = ""
        type        = "success"
      }
    })
  }
  trigger = {
    type = "EVENT"
    attributes = {
      id = "idn:identity-created"
    }
  }
}
