package v1alpha1

// fileHandler handler to manage file changes.
func (f *File) fileHandler() {
	// handle file removal
	if !f.Spec.Exists {
		f.fileRemoveHandler()
	}

	// handle resource modes
	if f.Spec.GetMode() != 0 {
		f.doFileMode()
	}

	return
}
