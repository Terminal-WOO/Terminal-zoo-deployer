package client

import "fmt"

// ErrInClusterConfig is a custom error type that wraps errors returned by
// rest.InClusterConfig from "k8s.io/client-go/rest". It provides additional
// context indicating a failure in loading the in-cluster Kubernetes configuration.
type ErrInClusterConfig struct {
	Err error // Err is the underlying error that occurred while building the in-cluster config.
}

// NewErrInClusterConfig creates a new ErrInClusterConfig error, wrapping the
// provided error returned from rest.InClusterConfig.
//
// Example usage:
//
//	config, err := rest.InClusterConfig()
//	if err != nil {
//	    return NewErrInClusterConfig(err)
//	}
func NewErrInClusterConfig(err error) *ErrInClusterConfig { return &ErrInClusterConfig{Err: err} }

// Error implements the error interface and returns a descriptive error message
// indicating that loading the in-cluster configuration failed.
func (e *ErrInClusterConfig) Error() string {
	return fmt.Sprintf("failed to load in-cluster config: %v", e.Err)
}

// Unwrap allows ErrInClusterConfig to support error unwrapping, enabling use of
// errors.Is and errors.As to inspect or match the underlying error.
func (e *ErrInClusterConfig) Unwrap() error { return e.Err }

// ErrOutClusterConfig is a custom error type that wraps an error returned when
// building an out-of-cluster Kubernetes configuration.
//
// It is typically used to wrap errors from clientcmd.BuildConfigFromFlags,
// which attempts to create a Kubernetes REST client configuration from explicit
// API server and kubeconfig path inputs.
//
// Example usage:
//
//	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
//	if err != nil {
//	    return NewErrOutClusterConfig(err)
//	}
type ErrOutClusterConfig struct {
	Err error // Err is the underlying error that occurred while building the out-of-cluster config.
}

// NewErrOutClusterConfig creates a new ErrOutClusterConfig by wrapping the provided error.
//
// This function is used to indicate that the failure occurred specifically during the
// construction of an out-of-cluster Kubernetes configuration.
func NewErrOutClusterConfig(err error) *ErrOutClusterConfig { return &ErrOutClusterConfig{Err: err} }

// Error implements the error interface for ErrOutClusterConfig.
//
// It returns a descriptive error message indicating the context of the failure.
func (e *ErrOutClusterConfig) Error() string {
	return fmt.Sprintf("failed to build out-of-cluster config: %v", e.Err)
}

// Unwrap allows unwrapping the underlying error.
//
// This supports Go's standard error unwrapping via errors.Unwrap and error inspection
// using errors.Is and errors.As.
func (e *ErrOutClusterConfig) Unwrap() error { return e.Err }

// ErrMissingEnvVars is returned when both the Unix and Windows home environment
// variables (typically "HOME" and "USERPROFILE") are not set in the environment.
//
// This error indicates that the application cannot determine the user's home directory
// because neither expected environment variable is available.
//
// You can check for this specific error type using errors.As or a type assertion.
type ErrMissingEnvVars struct{}

// Error implements the error interface for ErrMissingEnvVars.
func (e *ErrMissingEnvVars) Error() string {
	return fmt.Sprintf("required environment variables %q or %q are not set", UnixEnvHome, WindowsEnvHome)
}

// ErrClientCreation represents an error occurring during the creation of
// a Kubernetes clientset via kubernetes.NewForConfig.
//
// It wraps the underlying error returned by NewForConfig to provide context
// that the failure happened specifically during client creation.
//
// This error supports unwrapping to access the original error.
type ErrClientCreation struct {
	Err error // Err is the underlying error returned by kubernetes.NewForConfig.
}

// NewErrClientCreation constructs an ErrClientCreation that wraps the provided error.
//
//	err := errors.New("permission denied")
//	e := NewErrClientCreation(err)
//	fmt.Println(e.Error()) // prints: "failed to create Kubernetes client: permission denied"
func NewErrClientCreation(err error) *ErrClientCreation { return &ErrClientCreation{Err: err} }

// Error returns a descriptive error message including the wrapped error.
func (e *ErrClientCreation) Error() string {
	return fmt.Sprintf("failed to create Kubernetes client: %v", e.Err)
}

// Unwrap returns the underlying wrapped error.
func (e *ErrClientCreation) Unwrap() error { return e.Err }
