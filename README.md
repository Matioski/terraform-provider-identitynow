# IdentityNow Terraform Provider
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build And Test](https://github.com/SwissRe/terraform-provider-identitynow/actions/workflows/build-and-test-go.yaml/badge.svg)](https://github.com/SwissRe/terraform-provider-identitynow/actions/workflows/build-and-test-go.yaml)


Terraform provider is based on latest [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and [IdentityNow Go SDK](https://github.com/sailpoint-oss/golang-sdk)

## Usage
Include provider
```terraform
terraform {
  required_providers {
    identitynow = {
      source = "swissre/identitynow"
    }
  }
}
```

### Configure provider
For local usage you can use hardcoded provider configuration, but it is **highly** recommended to use environment variables.
```terraform
provider "identitynow" {
  host = "https://<tenant>.api.identitynow.com"
  client_id = "CLIENT_ID"
  client_secret = "CLIENT_SECRET"
}
```

Environment Variables:
* IDN_HOST
* IDN_CLIENT_ID
* IDN_CLIENT_SECRET

### Supported Terraform Data Sources
List of implemented data sources:
* Identity - `identitynow_identity`
* Cluster - `identitynow_cluster`
* Connector - `identitynow_connector`
* Entitlement - `identitynow_entitlement`

### Supported Terraform Resources
List of implemented resources:
* Identity Attribute - `identitynow_identity_attribute`
* Transform - `identitynow_transform`
* Source - `identitynow_source`
* Source Schema - `identitynow_source_schema`
* Identity Profile - `identitynow_identity_profile`
* Lifecycle State - `identitynow_lifecycle_state`
* Connector Rule - `identitynow_connector_rule`
* Workflow - `identitynow_workflow`

## Terraform Unstructured Object Type Workaround
Using `jsontype` as a workaround for unstructured data type
https://stackoverflow.com/questions/75024670/how-can-i-create-an-attribute-in-my-terraform-plugin-that-accepts-multiple-data


## Local Development
### Install Go
[Download](https://go.dev/doc/install) and install Go

### Install Terraform
[Download](https://www.terraform.io/downloads) latest terraform (>= 1.7.2) and add it to the PATH.

### Configure Terraform Plugin Framework
For local development follow up steps - https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install

In short: create `terraform.rc` in your `%APPDATA%` folder with similar content
```
provider_installation {

  dev_overrides {
      "swissre/identitynow" = "C:\\Users\\<USER_ID>\\go\\bin"
      }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

### Build and Install
Download dependencies
```shell
go mod download
```
Build provider
```shell
go install .
```

## Integration Tests
Terraform Plugin Framework supports integration tests - and it is recommended to use mock server for testing.
For mocking IdentityNow API we have used [Mockoon](https://mockoon.com/) which allows to create and run mock servers.

### Mockoon Configuration
1. [Download](https://mockoon.com/download/) and install Mockoon Desktop App (use portable version)
2. Import `mock/identitynow_mockoon.json` configuration file

### Run Mockoon from Command Line
Mockoon has CLI and Docker image. If possible use Docker image, but if you don't have Docker installed you can use CLI.

#### Mockoon Docker
```shell
docker run -d --mount type=bind,source=./mock/identitynow_mockoon.json,target=/data,readonly -p 3000:3000 mockoon/cli:latest -d data -p 3000
```

#### Mockoon CLI
```shell
mockoon-cli start --data ./mock/identitynow_mockoon.json --port 3000
```

### Run Acceptance Tests (using mock)
To enable Terraform Testing set environment variable `TF_ACC=1` and run tests
```shell
go test ./... -v
```

### Run Integration Tests
To enable Terraform Testing set environment variable `TF_ACC=1`, `IDN_HOST`, `IDN_CLIENT_ID` and `IDN_CLIENT_SECRET` and run integration tests
```shell
go test ./... -v -tags=integration
```

## Documentation
Documentation is generated using `tfplugindocs` tool. To generate documentation run
```shell
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
```


# Release
To release new version of the provider follow-up steps:
* From `main` branch search for latest tag and increment it by 1
```bash
git tag
git tag vX.Y.Z
```
* Push tag to the repository
```bash
git push origin vX.Y.Z
```
