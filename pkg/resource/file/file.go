package file

import (
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/state"
	"github.com/retr0h/psion/internal/status"
	"github.com/retr0h/psion/pkg/resource/api"
)

// GetExists the exists property.
func (f *Spec) GetExists() bool { return f.Exists }

// GetPath the path property.
func (f *Spec) GetPath() string { return f.Path }

// GetMode the mode property.
func (f *Spec) GetMode() fs.FileMode { return fs.FileMode(f.Mode) }

// GetModeString the mode property as an octal string.
func (f *Spec) GetModeString() string {
	if f.Mode == 0 {
		return ""
	}

	return fmt.Sprintf("0%o", f.Mode)
}

// New create a new instance of File.
func New(
	logger *slog.Logger,
	fileManager internal.FileManager,
	plan bool,
) *File {
	return &File{
		Status: &status.Status{
			Conditions: make([]status.StatusConditions, 0, 1),
		},
		Spec:   &Spec{},
		plan:   plan,
		logger: logger,
		file:   fileManager,
	}
}

// GetStatus determine the resources status.
func (f *File) GetStatus() api.Phase {
	noop := f.allMatch(api.NoOp)
	// set status to `NoOp` when all condition statuses are `NoOp`
	if noop {
		return api.NoOp
	}

	succeeded := f.allMatch(api.Succeeded)
	// set status to `Succeeded` when all condition statuses are `Succeeded`
	if succeeded {
		return api.Succeeded
	}

	pending := f.anyMatch(api.Pending)
	// set status to `Pending` when any condition statuses are `Pending`
	if pending {
		return api.Pending
	}

	failed := f.anyMatch(api.Failed)
	// set status to `Failed` when any condition statuses are `Failed`
	if failed {
		return api.Failed
	}

	// otherwise set to `Unknown`
	return api.Unknown
}

// GetStatusString determine the resources status as a string.
func (f *File) GetStatusString() string { return string(f.GetStatus()) }

func (f *File) allMatch(phase api.Phase) bool {
	for _, condition := range f.Status.Conditions {
		if condition.GetStatus() != phase {
			return false
		}
	}
	return true
}

func (f *File) anyMatch(phase api.Phase) bool {
	for _, condition := range f.Status.Conditions {
		if condition.GetStatus() == phase {
			return true
		}
	}
	return false
}

// GetStatusConditions the conditions property.
func (f *File) GetStatusConditions() []status.StatusConditions { return f.Status.Conditions }

// SetStatusCondition set the status condition property.
func (f *File) SetStatusCondition(
	statusType api.SpecAction,
	status_ api.Phase,
	message string,
	got string,
	want string,
) {
	reason := api.Apply
	if f.plan {
		reason = api.Plan
	}

	fileStatusConditions := status.StatusConditions{
		Type:    statusType,
		Status:  status_,
		Message: message,
		Reason:  reason,
		Got:     got,
		Want:    want,
	}

	conditions := f.Status.Conditions
	conditions = append(conditions, fileStatusConditions)
	f.Status.Conditions = conditions
}

// GetState provide current state after apply.
func (f *File) GetState() *state.StateResource {
	return &state.StateResource{
		Name:       f.Name,
		Kind:       f.Kind,
		APIVersion: f.APIVersion,
		Phase:      f.GetStatus(),
		Status: &status.Status{
			Phase:      f.GetStatus(),
			Conditions: f.GetStatusConditions(),
		},
	}
}

// Reconcile make consistent with the desired state.
func (f *File) Reconcile() error {
	f.fileHandler()

	f.logger.Info(
		"completed",
		slog.String("Status", f.GetStatusString()),
		slog.String("Kind", f.Kind),
		slog.String("APIVersion", f.APIVersion),
		slog.Group(Kind,
			slog.String("Path", f.Spec.GetPath()),
			slog.Bool("Exists", f.Spec.GetExists()),
			slog.String("Mode", f.Spec.GetModeString()),
		),
		slog.Group("Conditions", f.logStatusConditionGroups()...),
	)

	return nil
}

func (f *File) logStatusConditionGroups() []any {
	logGroups := make([]any, 0, len(f.Status.Conditions))
	for _, condition := range f.Status.Conditions {
		group := slog.Group(condition.GetTypeString(),
			slog.String("Status", condition.GetStatusString()),
			slog.String("Message", condition.GetMessage()),
			slog.String("Reason", condition.GetReasonString()),
			slog.String("Got", condition.GetGot()),
			slog.String("Want", condition.GetWant()),
		)

		logGroups = append(logGroups, group)
	}

	return logGroups
}
