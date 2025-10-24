package client

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	WindowsEnvHome string = "USERPROFILE"
	UnixEnvHome    string = "HOME"
	InClusterEnv   string = "KUBERNETES_SERVICE_HOST"
	KubeDir        string = ".kube"
	KubeConfig     string = "config"
)

var (
	// allows us to stub in tests
	inClusterConfigFunc    = rest.InClusterConfig
	outOfClusterConfigFunc = getOutOfClusterConfig
	getClusterConfigFunc   = getClusterConfig
	newForConfigFunc       = kubernetes.NewForConfig
)

// ConfigProvider defines a function that builds a config.
type ConfigProvider func(masterUrl, kubeConfigPath string) (*rest.Config, error)

// getHomeDir returns the user's home directory by checking the environment variables
// UnixEnvHome (e.g., "HOME") and WindowsEnvHome (e.g., "USERPROFILE") in that order.
//
// If neither environment variable is set, it returns an *ErrMissingEnvVars error.
//
// This function is useful for obtaining the home directory in a cross-platform way.
func getHomeDir() (string, error) {
	if home := os.Getenv(UnixEnvHome); home != "" {
		return home, nil
	}

	if home := os.Getenv(WindowsEnvHome); home != "" {
		return home, nil
	}

	return "", new(ErrMissingEnvVars)
}

// getOutOfClusterConfig returns the out-of-cluster configuration for accessing the Kubernetes API
// using the kubeconfig file located in the user's home directory.
// It is used when the code is running outside a Kubernetes cluster.
//
// The function attempts to locate the user's home directory to find the kubeconfig file.
// If the home directory cannot be determined from the environment, it returns an *ErrMissingEnvVars.
// If building the configuration fails, it returns an *ErrOutClusterConfig wrapping the original error.
//
// Returns:
//   - config: The configuration object for accessing the Kubernetes API using kubeconfig.
//   - error: An error object if the configuration cannot be created.
//
// Possible Errors:
//   - *ErrMissingEnvVars: Returned when neither the Unix nor Windows home environment variables are set.
//   - *ErrOutClusterConfig: Returned when the kubeconfig file cannot be loaded or parsed.
func getOutOfClusterConfig(buildConfig ConfigProvider) (*rest.Config, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		var expectedErr *ErrMissingEnvVars
		if !errors.As(err, &expectedErr) {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}

		return nil, err
	}

	config, err := buildConfig("", filepath.Join(homeDir, KubeDir, KubeConfig))
	if err != nil {
		return nil, NewErrOutClusterConfig(err)
	}

	return config, nil
}

// getClusterConfig returns the appropriate Kubernetes cluster configuration
// depending on whether the code is running inside or outside the cluster.
//
// When running inside a cluster, it uses the in-cluster configuration.
// When running outside a cluster, it falls back to out-of-cluster kubeconfig.
//
// Returns:
//   - config: The configuration object for accessing the Kubernetes API.
//   - error: An error object if the configuration cannot be created.
//
// Possible Errors:
//   - *ErrInClusterConfig: Returned if the in-cluster configuration cannot be loaded.
//   - *ErrMissingEnvVars: Returned if the home directory environment variables are missing when loading out-of-cluster config.
//   - *ErrOutClusterConfig: Returned if the out-of-cluster kubeconfig cannot be loaded or parsed.
func getClusterConfig() (*rest.Config, error) {
	if os.Getenv(InClusterEnv) != "" {
		config, err := inClusterConfigFunc()
		if err != nil {
			return nil, NewErrInClusterConfig(err)
		}

		return config, nil
	}

	config, err := outOfClusterConfigFunc(clientcmd.BuildConfigFromFlags)
	if err != nil {
		// Check for allowed error types
		var missingEnvErr *ErrMissingEnvVars
		var outClusterErr *ErrOutClusterConfig

		if !errors.As(err, &missingEnvErr) && !errors.As(err, &outClusterErr) {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}

		// Return the expected error types normally
		return nil, err
	}

	return config, nil
}

// CreateKubernetesClient creates a Kubernetes clientset based on the current environment configuration.
// It determines whether to use in-cluster or out-of-cluster configuration.
//
// Returns:
//   - clientset: The Kubernetes clientset for interacting with the Kubernetes API.
//   - error: An error object if the client cannot be created.
//
// Possible Errors:
//   - *ErrInClusterConfig: Returned if in-cluster configuration cannot be loaded.
//   - *ErrMissingEnvVars: Returned if required environment variables for out-of-cluster config are missing.
//   - *ErrOutClusterConfig: Returned if out-of-cluster kubeconfig cannot be loaded or parsed.
//   - *ErrClientCreation: Returned if the Kubernetes clientset cannot be created from the configuration.
func CreateKubernetesClient() (*kubernetes.Clientset, *rest.Config, error) {
	config, err := getClusterConfigFunc()
	if err != nil {
		// Check for allowed error types from getClusterConfig
		var missingEnvErr *ErrMissingEnvVars
		var outClusterErr *ErrOutClusterConfig
		var inClusterErr *ErrInClusterConfig

		if !errors.As(err, &missingEnvErr) &&
			!errors.As(err, &outClusterErr) &&
			!errors.As(err, &inClusterErr) {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}

		// Return expected errors normally
		return nil, nil, err
	}

	client, err := newForConfigFunc(config)
	if err != nil {
		return nil, nil, NewErrClientCreation(err)
	}

	return client, config, nil
}

// CreateControllerRuntimeClient creates a controller-runtime client.Client based on the current environment.
// It uses the global scheme, which must include any required CRDs (like Kueue).
//
// Returns:
//   - ctrlclient.Client: A generic Kubernetes client that supports CRUD for runtime.Objects.
//   - error: If the client cannot be created or configuration fails.
func CreateControllerRuntimeClient() (ctrlclient.Client, error) {
	config, err := getClusterConfigFunc()
	if err != nil {
		// Check for allowed error types
		var missingEnvErr *ErrMissingEnvVars
		var outClusterErr *ErrOutClusterConfig
		var inClusterErr *ErrInClusterConfig

		if !errors.As(err, &missingEnvErr) &&
			!errors.As(err, &outClusterErr) &&
			!errors.As(err, &inClusterErr) {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}

		return nil, err
	}

	cli, err := ctrlclient.New(config, ctrlclient.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create controller-runtime client: %w", err)
	}

	return cli, nil
}
