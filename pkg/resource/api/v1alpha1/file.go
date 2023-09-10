package v1alpha1

import (
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/spf13/afero"

	"github.com/retr0h/psion/pkg/resource/api"
)

// GetExists the exists property.
func (f *FileSpec) GetExists() bool { return f.Exists }

// GetPath the path property.
func (f *FileSpec) GetPath() string { return f.Path }

// GetMode the mode property.
func (f *FileSpec) GetMode() fs.FileMode { return fs.FileMode(f.Mode) }

// GetModeString the mode property as an octal string.
func (f *FileSpec) GetModeString() string {
	if f.Mode == 0 {
		return ""
	}

	return fmt.Sprintf("0%o", f.Mode)
}

// NewFile create a new instance of File.
func NewFile(
	logger *slog.Logger,
	appFs afero.Fs,
	plan bool,
) *File {
	return &File{
		Status: &api.Status{
			Conditions: make([]api.StatusConditions, 0, 1),
		},
		Spec:   &FileSpec{},
		plan:   plan,
		logger: logger,
		appFs:  appFs,
	}
}

// GetStatus determine the resources status.
func (f *File) GetStatus() api.Phase {
	noop := f.allConditionsMatch(api.NoOp, f.Status.Conditions)
	// set status to `NoOp` when all condition statuses are `NoOp`
	if noop {
		return api.NoOp
	}

	succeeded := f.allConditionsMatch(api.Succeeded, f.Status.Conditions)
	// set status to `Succeeded` when all condition statuses are `Succeeded`
	if succeeded {
		return api.Succeeded
	}

	pending := f.anyConditionsMatch(api.Pending, f.Status.Conditions)
	// set status to `Pending` when any condition statuses are `Pending`
	if pending {
		return api.Pending
	}

	failed := f.anyConditionsMatch(api.Failed, f.Status.Conditions)
	// set status to `Failed` when any condition statuses are `Failed`
	if failed {
		return api.Failed
	}

	// otherwise set to `Unknown`
	return api.Unknown
}

// GetStatusString determine the resources status as a string.
func (f *File) GetStatusString() string { return string(f.GetStatus()) }

func (f *File) allConditionsMatch(phase api.Phase, conditions []api.StatusConditions) bool {
	for _, condition := range conditions {
		if condition.GetStatus() != phase {
			return false
		}
	}
	return true
}

func (f *File) anyConditionsMatch(phase api.Phase, conditions []api.StatusConditions) bool {
	for _, condition := range conditions {
		if condition.GetStatus() == phase {
			return true
		}
	}
	return false
}

// GetStatusConditions the conditions property.
func (f *File) GetStatusConditions() []api.StatusConditions { return f.Status.Conditions }

// SetStatusCondition set the status condition property.
func (f *File) SetStatusCondition(
	statusType string,
	status api.Phase,
	message string,
	got string,
	want string,
) {
	reason := api.Apply
	if f.plan {
		reason = api.Plan
	}

	fileStatusConditions := api.StatusConditions{
		Type:    statusType,
		Status:  status,
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
func (f *File) GetState() *api.StateResource {
	return &api.StateResource{
		Name:       f.Name,
		Kind:       f.Kind,
		APIVersion: f.APIVersion,
		Phase:      f.GetStatus(),
		Status: &api.Status{
			Phase:      f.GetStatus(),
			Conditions: f.GetStatusConditions(),
		},
	}
}

// Reconcile make consistent with the desired state.
func (f *File) Reconcile() error {
	if !f.Spec.Exists {
		// file should be deleted
		f.fileRemoveHandler()
	} else {
		// handle remaining file permutations
		f.fileHandler()
	}

	f.logger.Info(
		"completed",
		slog.String("Status", f.GetStatusString()),
		slog.String("Kind", f.Kind),
		slog.String("APIVersion", f.APIVersion),
		slog.Group(FileKind,
			slog.String("Path", f.Spec.GetPath()),
			slog.Bool("Exists", f.Spec.GetExists()),
			slog.String("Mode", f.Spec.GetModeString()),
		),
		slog.Group("Conditions", f.logStatusConditionGroups()...),
	)

	return nil
}

func (f *File) logStatusConditionGroups() []any {
	var logGroups []any
	for _, condition := range f.Status.Conditions {
		group := slog.Group(condition.GetType(),
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
