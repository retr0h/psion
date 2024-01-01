package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config is a static copy of the resource's desired state.
type Config struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	metav1.TypeMeta   `json:",omitempty,inline"`
}
