package config

import (
	"embed"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/retr0h/psion/internal"
)

// Config is a static copy of the resource's desired state.
type Config struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	metav1.TypeMeta   `json:",omitempty,inline"`

	// file manager repository.
	fm internal.FileManager
	// embeded virtual file system.
	eFs embed.FS
}
