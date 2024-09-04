//go:build integration

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	"terraform-provider-identitynow/internal/util"
	"testing"
)

func TestIntegrationRole_CreateRoleAndUpdateMissingFields(t *testing.T) {
	entitlements := getEntitlements(2, context.Background())
	entitlementOne := entitlements[0]
	entitlementTwo := entitlements[1]

	accessProfileOne := getAccessProfiles(1, context.Background())[0]

	segment := getSegments(1, context.Background())[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndUpdateMissingFields"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndUpdateMissingFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", ownerIdentityName),
				),
			},
			// Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndUpdateMissingFields"
  description = "TestIntegrationRole_CreateRoleAndUpdateMissingFields"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name = "` + updatedOwnerIdentityName + `"
  }
  access_profiles = [
    {
      id   = "` + *accessProfileOne.Id + `"
      name = "` + accessProfileOne.Name + `"
    }
  ]
  entitlements = [
    {
        id = "` + *entitlementOne.Id + `"
        name = "` + *entitlementOne.Name + `"
    },
    {
        id = "` + *entitlementTwo.Id + `"
        name = "` + *entitlementTwo.Name + `"
    }
  ]
  requestable = true
  enabled = true
  access_request_config = {
    comments_required = true
    denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  revocation_request_config = {
    # Attributes defined in schema, but not used by IdentityNow (also not possible to set them from UI)
    #comments_required = true
    #denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  segments = [
    "` + *segment.Id + `"
  ]
 }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndUpdateMissingFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "description", "TestIntegrationRole_CreateRoleAndUpdateMissingFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", updatedOwnerIdentityName),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "access_profiles.*", map[string]string{
						"type": "ACCESS_PROFILE",
						"id":   *accessProfileOne.Id,
						"name": accessProfileOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementOne.Id,
						"name": *entitlementOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementTwo.Id,
						"name": *entitlementTwo.Name,
					}),
					resource.TestCheckResourceAttr("identitynow_role.test", "requestable", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.1.approver_type", "MANAGER"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.comments_required", "true"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.1.approver_type", "MANAGER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "segments.0", *segment.Id),
					//resource.TestCheckResourceAttr("identitynow_role.test", "legacy_membership_info", "STANDARD"),
				),
			},
		},
	})
}

func TestIntegrationRole_CreateRoleAndRemoveOptionalFields(t *testing.T) {
	entitlements := getEntitlements(2, context.Background())
	entitlementOne := entitlements[0]
	entitlementTwo := entitlements[1]

	accessProfileOne := getAccessProfiles(1, context.Background())[0]
	segment := getSegments(1, context.Background())[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndUpdateMissingFields"
  description = "TestIntegrationRole_CreateRoleAndUpdateMissingFields"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  access_profiles = [
    {
      id   = "` + *accessProfileOne.Id + `"
      name = "` + accessProfileOne.Name + `"
    }
  ]
  entitlements = [
    {
        id = "` + *entitlementOne.Id + `"
        name = "` + *entitlementOne.Name + `"
    },
    {
        id = "` + *entitlementTwo.Id + `"
        name = "` + *entitlementTwo.Name + `"
    }
  ]
  requestable = true
  enabled = true
  access_request_config = {
    comments_required = true
    denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  revocation_request_config = {
    # Attributes defined in schema, but not used by IdentityNow (also not possible to set them from UI)
    #comments_required = true
    #denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  segments = [
    "` + *segment.Id + `"
  ]
 }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndUpdateMissingFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "description", "TestIntegrationRole_CreateRoleAndUpdateMissingFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", ownerIdentityName),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "access_profiles.*", map[string]string{
						"type": "ACCESS_PROFILE",
						"id":   *accessProfileOne.Id,
						"name": accessProfileOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementOne.Id,
						"name": *entitlementOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementTwo.Id,
						"name": *entitlementTwo.Name,
					}),
					resource.TestCheckResourceAttr("identitynow_role.test", "requestable", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.1.approver_type", "MANAGER"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.comments_required", "true"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.1.approver_type", "MANAGER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "segments.0", *segment.Id),
					//resource.TestCheckResourceAttr("identitynow_role.test", "legacy_membership_info", "STANDARD"),
				),
			},
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndUpdateMissingFields"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndUpdateMissingFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", ownerIdentityName),
					resource.TestCheckNoResourceAttr("identitynow_role.test", "access_profiles"),
					resource.TestCheckNoResourceAttr("identitynow_role.test", "revocation_request_config"),
				),
			},
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndUpdateMissingFields"
  description = "TestIntegrationRole_CreateRoleAndUpdateMissingFields"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  access_profiles = [
    {
      id   = "` + *accessProfileOne.Id + `"
      name = "` + accessProfileOne.Name + `"
    }
  ]
  entitlements = [
    {
        id = "` + *entitlementOne.Id + `"
        name = "` + *entitlementOne.Name + `"
    },
    {
        id = "` + *entitlementTwo.Id + `"
        name = "` + *entitlementTwo.Name + `"
    }
  ]
  requestable = true
  enabled = true
  access_request_config = {
    comments_required = true
    denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  revocation_request_config = {
    # Attributes defined in schema, but not used by IdentityNow (also not possible to set them from UI)
    #comments_required = true
    #denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  segments = [
    "` + *segment.Id + `"
  ]
 }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndUpdateMissingFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "description", "TestIntegrationRole_CreateRoleAndUpdateMissingFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", ownerIdentityName),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "access_profiles.*", map[string]string{
						"type": "ACCESS_PROFILE",
						"id":   *accessProfileOne.Id,
						"name": accessProfileOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementOne.Id,
						"name": *entitlementOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementTwo.Id,
						"name": *entitlementTwo.Name,
					}),
					resource.TestCheckResourceAttr("identitynow_role.test", "requestable", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "enabled", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.1.approver_type", "MANAGER"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.comments_required", "true"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.1.approver_type", "MANAGER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "segments.0", *segment.Id),
					//resource.TestCheckResourceAttr("identitynow_role.test", "legacy_membership_info", "STANDARD"),
				),
			},
		},
	})
}

func TestIntegrationRole_CreateRoleAndEditFields(t *testing.T) {
	entitlements := getEntitlements(2, context.Background())
	entitlementOne := entitlements[0]
	entitlementTwo := entitlements[1]

	accessProfiles := getAccessProfiles(2, context.Background())
	accessProfileOne := accessProfiles[0]
	accessProfileTwo := accessProfiles[1]

	segments := getSegments(2, context.Background())
	segmentOne := segments[0]
	segmentTwo := segments[1]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndEditFields"
  description = "TestIntegrationRole_CreateRoleAndEditFields"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  access_profiles = [
    {
      id   = "` + *accessProfileOne.Id + `"
      name = "` + accessProfileOne.Name + `"
    }
  ]
  entitlements = [
    {
        id = "` + *entitlementOne.Id + `"
        name = "` + *entitlementOne.Name + `"
    }
  ]
  requestable = false
  enabled = false
  access_request_config = {
    comments_required = true
    denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  revocation_request_config = {
    # Attributes defined in schema, but not used by IdentityNow (also not possible to set them from UI)
    #comments_required = true
    #denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  segments = [
    "` + *segmentOne.Id + `"
  ]
 }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndEditFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "description", "TestIntegrationRole_CreateRoleAndEditFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", ownerIdentityName),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "access_profiles.*", map[string]string{
						"type": "ACCESS_PROFILE",
						"id":   *accessProfileOne.Id,
						"name": accessProfileOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementOne.Id,
						"name": *entitlementOne.Name,
					}),
					resource.TestCheckResourceAttr("identitynow_role.test", "requestable", "false"),
					resource.TestCheckResourceAttr("identitynow_role.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.1.approver_type", "MANAGER"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.comments_required", "true"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.1.approver_type", "MANAGER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "segments.0", *segmentOne.Id),
					//resource.TestCheckResourceAttr("identitynow_role.test", "legacy_membership_info", "STANDARD"),
				),
			}, // Update and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndEditFields"
  description = "TestIntegrationRole_CreateRoleAndEditFields update"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name = "` + updatedOwnerIdentityName + `"
  }
  access_profiles = [
    {
      id   = "` + *accessProfileOne.Id + `"
      name = "` + accessProfileOne.Name + `"
    },
    {
      id   = "` + *accessProfileTwo.Id + `"
      name = "` + accessProfileTwo.Name + `"
    }
  ]
  entitlements = [
    {
        id = "` + *entitlementOne.Id + `"
        name = "` + *entitlementOne.Name + `"
    },
    {
        id = "` + *entitlementTwo.Id + `"
        name = "` + *entitlementTwo.Name + `"
    }
  ]
  requestable = false
  enabled = false
  access_request_config = {
    comments_required = false
    denial_comments_required = false
    approval_schemas = [
        {
            approver_type = "OWNER"
        }
    ]
  }
  revocation_request_config = {
    # Attributes defined in schema, but not used by IdentityNow (also not possible to set them from UI)
    #comments_required = true
    #denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        }
    ]
  }
  segments = [
    "` + *segmentOne.Id + `",
    "` + *segmentTwo.Id + `"
  ]
 }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndEditFields"),
					resource.TestCheckResourceAttr("identitynow_role.test", "description", "TestIntegrationRole_CreateRoleAndEditFields update"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", updatedOwnerIdentityName),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "access_profiles.*", map[string]string{
						"type": "ACCESS_PROFILE",
						"id":   *accessProfileOne.Id,
						"name": accessProfileOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "access_profiles.*", map[string]string{
						"type": "ACCESS_PROFILE",
						"id":   *accessProfileTwo.Id,
						"name": accessProfileTwo.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementOne.Id,
						"name": *entitlementOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementTwo.Id,
						"name": *entitlementTwo.Name,
					}),
					resource.TestCheckResourceAttr("identitynow_role.test", "requestable", "false"),
					resource.TestCheckResourceAttr("identitynow_role.test", "enabled", "false"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.comments_required", "false"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.denial_comments_required", "false"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.0.approver_type", "OWNER"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.comments_required", "true"),
					//resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "segments.0", *segmentOne.Id),
					resource.TestCheckResourceAttr("identitynow_role.test", "segments.1", *segmentTwo.Id),
					//resource.TestCheckResourceAttr("identitynow_role.test", "legacy_membership_info", "STANDARD"),
				),
			},
		},
	})
}

func TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation(t *testing.T) {
	entitlements := getEntitlements(2, context.Background())
	entitlementOne := entitlements[0]
	entitlementTwo := entitlements[1]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation"
  description = "TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation"
  owner = {
    id   = "` + ownerIdentityId + `"
    name = "` + ownerIdentityName + `"
  }
  entitlements = [
    {
        id = "` + *entitlementOne.Id + `"
        name = "` + *entitlementOne.Name + `"
    }
  ]
  requestable = true
  enabled = true
  access_request_config = {
    comments_required = true
    denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        }
    ]
  }
  revocation_request_config = {
    # Attributes defined in schema, but not used by IdentityNow (also not possible to set them from UI)
    #comments_required = true
    #denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        }
    ]
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation"),
					resource.TestCheckResourceAttr("identitynow_role.test", "description", "TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", ownerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", ownerIdentityName),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementOne.Id,
						"name": *entitlementOne.Name,
					}),
					resource.TestCheckResourceAttr("identitynow_role.test", "requestable", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.0.approver_type", "OWNER"),
				),
			}, // Create and Read testing
			{
				Config: providerIntegrationConfig + `
resource "identitynow_role" "test" {
  name        = "TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation updated"
  description = "TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation updated"
  owner = {
    id   = "` + updatedOwnerIdentityId + `"
    name = "` + updatedOwnerIdentityName + `"
  }
  entitlements = [
    {
        id = "` + *entitlementOne.Id + `"
        name = "` + *entitlementOne.Name + `"
    },
    {
        id = "` + *entitlementTwo.Id + `"
        name = "` + *entitlementTwo.Name + `"
    }
  ]
  requestable = true
  enabled = true
  access_request_config = {
    comments_required = true
    denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
  revocation_request_config = {
    # Attributes defined in schema, but not used by IdentityNow (also not possible to set them from UI)
    #comments_required = true
    #denial_comments_required = true
    approval_schemas = [
        {
            approver_type = "OWNER"
        },
        {
            approver_type = "MANAGER"
        }
    ]
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("identitynow_role.test", "id"),
					resource.TestCheckResourceAttr("identitynow_role.test", "name", "TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation updated"),
					resource.TestCheckResourceAttr("identitynow_role.test", "description", "TestIntegrationRole_CreateRoleAndEditFields_TriggersRecreation updated"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.type", "IDENTITY"),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.id", updatedOwnerIdentityId),
					resource.TestCheckResourceAttr("identitynow_role.test", "owner.name", updatedOwnerIdentityName),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementOne.Id,
						"name": *entitlementOne.Name,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("identitynow_role.test", "entitlements.*", map[string]string{
						"type": "ENTITLEMENT",
						"id":   *entitlementTwo.Id,
						"name": *entitlementTwo.Name,
					}),
					resource.TestCheckResourceAttr("identitynow_role.test", "requestable", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.denial_comments_required", "true"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "access_request_config.approval_schemas.1.approver_type", "MANAGER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.0.approver_type", "OWNER"),
					resource.TestCheckResourceAttr("identitynow_role.test", "revocation_request_config.approval_schemas.1.approver_type", "MANAGER"),
				),
			}},
	})
}

func getEntitlements(limit int32, ctx context.Context) []sailpointBeta.Entitlement {
	entitlements, spResp, err := SPApiClient.Beta.EntitlementsAPI.ListEntitlements(context.Background()).Limit(limit).Execute()
	if err != nil {
		fmt.Printf("Error fetching entitlements: %s\n%s", err, util.GetBody(spResp))
	}
	if len(entitlements) < int(limit) {
		fmt.Printf("Unable to provide %d entitlements; %d found\n", limit, len(entitlements))
	}
	return entitlements
}

func getSegments(limit int32, ctx context.Context) []sailpointBeta.Segment {
	segments, spResp, err := SPApiClient.Beta.SegmentsAPI.ListSegments(context.Background()).Limit(limit).Execute()
	if err != nil {
		fmt.Printf("Error fetching segments: %s\n%s", err, util.GetBody(spResp))
	}
	if len(segments) < int(limit) {
		fmt.Printf("Unable to provide %d segments; %d found\n", limit, len(segments))
	}
	return segments
}
