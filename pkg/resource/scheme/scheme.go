package scheme

import (
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Scheme struct {
	// APIVersion version of API resource to use.
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	// Kind is the type of resource being referenced.
	Kind string `json:"kind" yaml:"kind"`

	logger *logrus.Logger
}

// Scheme create a new instance of Scheme.
func New(
	logger *logrus.Logger,
) *Scheme {
	return &Scheme{
		logger: logger,
	}
}

// Decode YAML resource into a scheme to be used.
func (s *Scheme) Decode(
	resource []byte,
) *Scheme {
	var resourceData *Scheme
	err := yaml.Unmarshal(resource, &resourceData)
	if err != nil {
		logrus.WithError(err).
			Fatal("cannot unmarshal")
	}

	return resourceData
}

// Version of the API resource.
func (s *Scheme) Version() string {
	parts := strings.Split(s.APIVersion, "/")

	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

// switch resourceData["apiVersion"] {
// case "files.psion.io/v1alpha1":
// 	var resourceKind v1alpha1.Manager
// 	switch resourceData["kind"] {
// 	case "File":
// 		resourceKind = v1alpha1.NewFile(
// 			logger,
// 			appFs,
// 			resourceData,
// 		)
// 		resourceKind.Reconcile()
// 	default:
// 		logger.WithFields(logrus.Fields{
// 			"kind": resourceData["kind"],
// 		}).Fatal("unsupported kind")
// 	}
// default:
// 	logger.WithFields(logrus.Fields{
// 		"apiVersion": resourceData["apiVersion"],
// 	}).Fatal("unsupported apiVersion")
// }
