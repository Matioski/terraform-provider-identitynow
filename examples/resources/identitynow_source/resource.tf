resource "identitynow_source" "demo_source" {
  name        = "DEMO Terraform IDN Loopback"
  description = "Creating from terraform"
  owner = {
    id   = data.identitynow_identity.default_owner.id
    name = data.identitynow_identity.default_owner.name
  }
  cluster = {
    id   = "127318236adsasdjasdas123"
    name = "sp_connect_proxy_cluster"
  }
  features = []
  connector_attributes = jsonencode({
    enableLCS           = true
    test_attribute_1    = "test-1 change"
    anotherBigAttribute = "test-2 change"
    inherited = {
      first   = "1"
      seconds = "2"
    }
  })
  connector_attributes_credentials = jsonencode({
    username = "test"
    password = "test"
  })
  type             = data.identitynow_connector.idn_connector.type
  connector        = data.identitynow_connector.idn_connector.type
  delete_threshold = 10
}
