# GCP Scheduled GCE Instance Start - PubSub Example

This example demonstrates how to use the
[root module][root-module] and the
[event-project-log-entry submodule][event-project-log-entry-submodule]
to configure a system
which responds to Compute VM creation events by labelling them with the
principal email address of the account responsible for causing the events.

## Usage

To provision this example, populate `terraform.tfvars` with the [required variables](#inputs) and run the following commands within
this directory:

- `terraform init` to initialize the directory
- `terraform plan` to generate the execution plan
- `terraform apply` to apply the execution plan
- `terraform destroy` to destroy the infrastructure

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| action | Instance action parameters | `map(any)` | <pre>{<br>  "start": {<br>    "query": "labels.start_daily:true AND state:TERMINATED",<br>    "schedule": "0 1 * * *"<br>  },<br>  "stop": {<br>    "query": "labels.stop_daily:true AND state:RUNNING",<br>    "schedule": "0 2 * * *"<br>  }<br>}</pre> | no |
| org\_id | The organization ID to which resources will be applied. | `string` | `"override in terraform.tfvars"` | no |
| project\_id | The ID of the project to which resources will be applied. | `string` | n/a | yes |
| region | The region in which resources will be applied. | `string` | n/a | yes |
| search\_scope | The scope of the search | `string` | `"projects"` | no |
| service\_account\_email | The service account email | `string` | `""` | no |
| time\_zone | The timezone to use in scheduler | `string` | `"Etc/UTC"` | no |
| vm | VM spec - zone and subnetwork. Null to disable | <pre>object({<br>    zone       = string<br>    subnetwork = string<br>  })</pre> | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| function\_name | The name of the function created |
| project\_id | The project in which resources are applied. |
| region | The region in which resources are applied. |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

## Requirements

The following sections describe the requirements which must be met in
order to invoke this module. The requirements of the
[root module][root-module-requirements] and the
[event-project-log-entry submodule][event-project-log-entry-submodule-requirements]
must also be met.

### Software Dependencies

The following software dependencies must be installed on the system
from which this module will be invoked:

- [Terraform][terraform-site] v0.12

### IAM Roles

The Service Account which will be used to invoke this module must have
the following IAM roles:

- Compute Instance Admin (v1): `roles/compute.instanceAdmin.v1`

### APIs

The project against which this module will be invoked must have the
following APIs enabled:

- Compute Engine API: `compute.googleapis.com`

[event-project-log-entry-submodule-requirements]: ../../modules/event-project-log-entry/README.md#requirements
[event-project-log-entry-submodule]: ../../modules/event-project-log-entry
[root-module-requirements]: ../../README.md#requirements
[root-module]: ../..
[terraform-site]: https://terraform.io/
