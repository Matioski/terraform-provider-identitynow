## [1.0.0] (October 03, 2024)
Initial version of IdentityNow Terraform Provider


### Added

* Provider Data Sources:
  * `identitynow_identity` - retrieve Identity by alias (username)
  * `identitynow_cluster` - retrieve Cluster by name
  * `identitynow_connector` - retrieve Connector by name
  * `identitynow_entitlement` - retrieve Entitlement by Source ID and Entitlement Value
* Provider Resources:
  * `identitynow_connector_rule` - manage Connector Rule
  * `identitynow_identity_attribute` - manage Identity Attribute
  * `identitynow_identity_profile` - manage Identity Profile
  * `identitynow_lifecycle_state` - manage Lifecycle State
  * `identitynow_source` - manage Source
  * `identitynow_source_schema` - manage Source Schema
  * `identitynow_transform` - manage Transform
  * `identitynow_workflow` - manage Workflow
