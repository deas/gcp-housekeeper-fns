locals {
  name = "audit-label"
}

resource "random_pet" "main" {
  separator = "-"
}

module "vm" {
  source     = "../../examples/vm"
  project_id = var.project_id
  zone       = var.vm.zone
  subnetwork = var.vm.subnetwork
  count      = var.vm != null ? 1 : 0
  depends_on = [module.audit_label]
}

# TODO: Delegates to
# https://github.com/terraform-google-modules/terraform-google-log-export
# Event function submodules also support folders - but not organization
module "event_log_entry" {
  source     = "../../modules/event-organization-log-entry"
  filter     = "protoPayload.@type=\"type.googleapis.com/google.cloud.audit.AuditLog\" protoPayload.methodName:insert operation.first=true"
  name       = "${local.name}-${random_pet.main.id}"
  project_id = var.project_id
  org_id     = var.org_id
}

/*
module "event_log_entry" {
  source     = "terraform-google-modules/event-function/google//modules/event-project-log-entry"
  version    = "2.3.0"
  filter     = "protoPayload.@type=\"type.googleapis.com/google.cloud.audit.AuditLog\" protoPayload.methodName:insert operation.first=true"
  name       = "${local.name}-${random_pet.main.id}"
  project_id = var.project_id
}
*/

# TODO: A bit hacky? Open for ideas.
module "function" {
  source = "../.."
}

module "audit_label" {
  source      = "terraform-google-modules/event-function/google"
  version     = "2.3.0"
  description = "Labels resource with owner information."
  entry_point = module.function.entry_point_v1["label"]
  #environment_variables = {
  #  LABEL_KEY = "principal-email"
  #}
  event_trigger                  = module.event_log_entry.function_event_trigger
  name                           = "${local.name}-${random_pet.main.id}"
  project_id                     = var.project_id
  region                         = var.region
  source_directory               = module.function.path
  files_to_exclude_in_source_dir = module.function.excludes
  available_memory_mb            = "128"
  runtime                        = module.function.runtime
}

/* Hopefully no longer needed
resource "null_resource" "wait_for_function" {
  provisioner "local-exec" {
    command = "sleep 60"
  }

  depends_on = [module.audit_label]
}
*/