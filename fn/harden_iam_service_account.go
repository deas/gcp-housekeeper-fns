package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func HardenPubSub(ctx context.Context, m PubSubMessage) error {
	Logger.Info(ctx, fmt.Sprintf("Got PubSub message %s", string(m.Data))) // Automatically decoded from base64
	logentry := &AuditLogEntry{}
	// logentry := &audit.AuditLog{}
	// var auditLogEntry AuditLogEntry
	err := json.Unmarshal(m.Data, &logentry)
	if err != nil {
		Logger.Info(ctx, fmt.Sprintf("Error: could not unmarshall to audit log %v\n", err))
	}
	return harden(ctx, *logentry.ProtoPayload, fmt.Sprintf("%s/%s", logentry.ProtoPayload.ServiceName, logentry.ProtoPayload.ResourceName))
}

func HardenEvent(ctx context.Context, ev event.Event) error {
	Logger.Info(ctx, fmt.Sprintf("Got CloudEvent %s with data %s", ev.ID(), string(ev.Data())))
	logentry := &AuditLogEntry{}
	if err := ev.DataAs(logentry); err != nil {
		return fmt.Errorf("Error parsing event payload : %w", err)
	}
	return harden(ctx, *logentry.ProtoPayload, ev.Subject())
}

func harden(ctx context.Context, payload AuditLogProtoPayload, subject string) error {
	/*
		// TODO: Add your project ID
		projectID := "your-project-id" // flag.String("project_id", "", "Cloud Project ID")
		// TODO: Add the ID of your member in the form "user:member@example.com"
		member := "your-memeber-id" // flag.String("member_id", "", "Your member ID")
		// flag.Parse()

		// The role to be granted
		var role string = "roles/logging.logWriter"

		// Initializes the Cloud Resource Manager service
		crmService, err := cloudresourcemanager.NewService(ctx)
		if err != nil {
			stdlog.Fatalf("cloudresourcemanager.NewService: %v", err)
		}

		// Grants your member the "Log writer" role for your project
		addBinding(ctx, crmService, projectID, member, role)

		// Gets the project's policy and prints all members with the "Log Writer" role
		policy := getPolicy(crmService, projectID)
		// Find the policy binding for role. Only one binding can have the role.
		var binding *cloudresourcemanager.Binding
		for _, b := range policy.Bindings {
			if b.Role == role {
				binding = b
				break
			}
		}
		log.Info(ctx, fmt.Sprintf("Role: %s", binding.Role))
		log.Info(ctx, fmt.Sprintf("Members: %s", strings.Join(binding.Members, ", ")))

		// Removes member from the "Log writer" role
		removeMember(ctx, crmService, projectID, member, role)
	*/
	return nil

}

// addBinding adds the member to the project's IAM policy
func addBinding(ctx context.Context, crmService *cloudresourcemanager.Service, projectID, member, role string) {

	policy := getPolicy(crmService, projectID)

	// Find the policy binding for role. Only one binding can have the role.
	var binding *cloudresourcemanager.Binding
	for _, b := range policy.Bindings {
		if b.Role == role {
			binding = b
			break
		}
	}

	if binding != nil {
		// If the binding exists, adds the member to the binding
		binding.Members = append(binding.Members, member)
	} else {
		// If the binding does not exist, adds a new binding to the policy
		binding = &cloudresourcemanager.Binding{
			Role:    role,
			Members: []string{member},
		}
		policy.Bindings = append(policy.Bindings, binding)
	}

	setPolicy(ctx, crmService, projectID, policy)

}

// removeMember removes the member from the project's IAM policy
func removeMember(ctx context.Context, crmService *cloudresourcemanager.Service, projectID, member, role string) {

	policy := getPolicy(crmService, projectID)

	// Find the policy binding for role. Only one binding can have the role.
	var binding *cloudresourcemanager.Binding
	var bindingIndex int
	for i, b := range policy.Bindings {
		if b.Role == role {
			binding = b
			bindingIndex = i
			break
		}
	}

	// Order doesn't matter for bindings or members, so to remove, move the last item
	// into the removed spot and shrink the slice.
	if len(binding.Members) == 1 {
		// If the member is the only member in the binding, removes the binding
		last := len(policy.Bindings) - 1
		policy.Bindings[bindingIndex] = policy.Bindings[last]
		policy.Bindings = policy.Bindings[:last]
	} else {
		// If there is more than one member in the binding, removes the member
		var memberIndex int
		for i, mm := range binding.Members {
			if mm == member {
				memberIndex = i
			}
		}
		last := len(policy.Bindings[bindingIndex].Members) - 1
		binding.Members[memberIndex] = binding.Members[last]
		binding.Members = binding.Members[:last]
	}

	setPolicy(ctx, crmService, projectID, policy)

}

// getPolicy gets the project's IAM policy
func getPolicy(crmService *cloudresourcemanager.Service, projectID string) *cloudresourcemanager.Policy {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	request := new(cloudresourcemanager.GetIamPolicyRequest)
	policy, err := crmService.Projects.GetIamPolicy(projectID, request).Do()
	if err != nil {
		log.Fatalf("Projects.GetIamPolicy: %v", err)
	}

	return policy
}

// setPolicy sets the project's IAM policy
func setPolicy(ctx context.Context, crmService *cloudresourcemanager.Service, projectID string, policy *cloudresourcemanager.Policy) {
	// ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	// defer cancel()
	request := new(cloudresourcemanager.SetIamPolicyRequest)
	request.Policy = policy
	policy, err := crmService.Projects.SetIamPolicy(projectID, request).Do()
	if err != nil {
		log.Fatalf("Projects.SetIamPolicy: %v", err)
	}
}
