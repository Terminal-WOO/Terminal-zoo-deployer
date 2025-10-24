package client

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	ConfigHost      string = "https://mock-k8s"
	UnixMockHome    string = "/home/mock"
	WindowsMockHome string = "C:\\Users\\MockUser"
)

// --- Helpers to stub out dependencies ---
func resetStubs() {
	getClusterConfigFunc = getClusterConfig
	newForConfigFunc = kubernetes.NewForConfig
}

func TestGetHomeDir(t *testing.T) {
	tests := []struct {
		name           string
		setUnixEnv     string
		setWindowsEnv  string
		wantDir        string
		wantErrMissing bool
	}{
		{
			name:           "Neither env set → ErrMissingEnvVars",
			setUnixEnv:     "",
			setWindowsEnv:  "",
			wantDir:        "",
			wantErrMissing: true,
		},
		{
			name:           "Only HOME set → return HOME",
			setUnixEnv:     UnixMockHome,
			setWindowsEnv:  "",
			wantDir:        UnixMockHome,
			wantErrMissing: false,
		},
		{
			name:           "Only USERPROFILE set → return USERPROFILE",
			setUnixEnv:     "",
			setWindowsEnv:  WindowsMockHome,
			wantDir:        WindowsMockHome,
			wantErrMissing: false,
		},
		{
			name:           "Both set → HOME takes precedence",
			setUnixEnv:     UnixMockHome,
			setWindowsEnv:  WindowsMockHome,
			wantDir:        UnixMockHome,
			wantErrMissing: false,
		},
	}

	for _, tc := range tests {
		// capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Use t.Setenv (Go 1.17+) to set and auto-restore.
			t.Setenv(UnixEnvHome, tc.setUnixEnv)
			t.Setenv(WindowsEnvHome, tc.setWindowsEnv)

			gotDir, err := getHomeDir()

			if tc.wantErrMissing {
				// 1) Expect a non-nil error
				if err == nil {
					t.Fatalf("expected ErrMissingEnvVars, got nil error and dir=%q", gotDir)
				}

				// 2) Assert error is *ErrMissingEnvVars via errors.As
				var missingEnvErr *ErrMissingEnvVars
				if !errors.As(err, &missingEnvErr) {
					t.Fatalf("expected error of type *ErrMissingEnvVars via errors.As, got %T: %v", err, err)
				}

				// 3) Explicit type assertion check
				if _, ok := err.(*ErrMissingEnvVars); !ok {
					t.Fatalf("expected error to be *ErrMissingEnvVars via type assertion, got %T", err)
				}

				// 4) Ensure returned dir is empty
				if gotDir != "" {
					t.Errorf("when error is ErrMissingEnvVars, expected empty dir, got %q", gotDir)
				}

			} else {
				// Expect no error and correct dir
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if gotDir != tc.wantDir {
					t.Errorf("gotDir = %q; want %q", gotDir, tc.wantDir)
				}
			}
		})
	}
}

func TestGetOutOfClusterConfig(t *testing.T) {
	t.Run("Missing home", func(t *testing.T) {
		// Ensure neither env var is set so getHomeDir returns *ErrMissingEnvVars
		t.Setenv(UnixEnvHome, "")
		t.Setenv(WindowsEnvHome, "")

		// buildConfig should never be called in this case
		badBuild := func(masterURL, kubeconfigPath string) (*rest.Config, error) {
			t.Fatal("buildConfig was called even though home dir lookup failed")
			return nil, nil
		}

		cfg, err := getOutOfClusterConfig(badBuild)
		if cfg != nil {
			t.Errorf("expected cfg to be nil on missing home, got %+v", cfg)
		}
		if err == nil {
			t.Fatalf("expected an error when both %s and %s are unset", UnixEnvHome, WindowsEnvHome)
		}

		// 1) Check via errors.As
		var missErr *ErrMissingEnvVars
		if !errors.As(err, &missErr) {
			t.Fatalf("expected error of type *ErrMissingEnvVars via errors.As, got %T: %v", err, err)
		}

		// 2) Direct type assertion
		if _, ok := err.(*ErrMissingEnvVars); !ok {
			t.Fatalf("expected error to be *ErrMissingEnvVars via type assertion, got %T", err)
		}
	})

	t.Run("Build success", func(t *testing.T) {
		// Set HOME so getHomeDir succeeds
		t.Setenv(UnixEnvHome, UnixMockHome)
		t.Setenv(WindowsEnvHome, "")

		// Capture the args passed into buildConfig and return a dummy rest.Config
		var gotURL, gotPath string
		fakeConfig := &rest.Config{Host: ConfigHost}
		buildSuccess := func(masterURL, kubeconfigPath string) (*rest.Config, error) {
			gotURL, gotPath = masterURL, kubeconfigPath
			return fakeConfig, nil
		}

		cfg, err := getOutOfClusterConfig(buildSuccess)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg != fakeConfig {
			t.Errorf("expected returned config %#v, got %#v", fakeConfig, cfg)
		}

		// masterURL must be empty, and path must be homeDir/KubeDir/KubeConfig
		expectedPath := filepath.Join(UnixMockHome, KubeDir, KubeConfig)
		if gotURL != "" {
			t.Errorf("expected masterURL == \"\", got %q", gotURL)
		}

		if gotPath != expectedPath {
			t.Errorf("expected kubeconfigPath == %q, got %q", expectedPath, gotPath)
		}
	})

	t.Run("Build failure", func(t *testing.T) {
		// Set USERPROFILE so getHomeDir picks WindowsEnvHome
		t.Setenv(UnixEnvHome, "")
		t.Setenv(WindowsEnvHome, WindowsMockHome)

		// Simulate buildConfig failure
		origErr := fmt.Errorf("parse error")
		buildFail := func(masterURL, kubeconfigPath string) (*rest.Config, error) {
			return nil, origErr
		}

		cfg, err := getOutOfClusterConfig(buildFail)
		if cfg != nil {
			t.Errorf("expected cfg to be nil on build failure, got %+v", cfg)
		}
		if err == nil {
			t.Fatal("expected an error when buildConfig fails")
		}

		// 1) Check via errors.As
		var outErr *ErrOutClusterConfig
		if !errors.As(err, &outErr) {
			t.Fatalf("expected error of type *ErrOutClusterConfig via errors.As, got %T: %v", err, err)
		}

		// 2) Direct type assertion
		if _, ok := err.(*ErrOutClusterConfig); !ok {
			t.Fatalf("expected error to be *ErrOutClusterConfig via type assertion, got %T", err)
		}

		// 3) Unwrap yields original error
		if !errors.Is(err, origErr) {
			t.Errorf("expected wrapped error to contain original; errors.Is(err, origErr) = false")
		}
	})
}

func TestGetClusterConfig(t *testing.T) {
	t.Run("In cluster error", func(t *testing.T) {
		// Simulate “inside cluster” by setting the marker env var.
		t.Setenv(InClusterEnv, "true")

		cfg, err := getClusterConfig()
		if cfg != nil {
			t.Fatalf("expected nil config on in-cluster failure, got %+v", cfg)
		}
		if err == nil {
			t.Fatal("expected an error when InClusterConfig fails")
		}

		// It should be wrapped in *ErrInClusterConfig
		var inErr *ErrInClusterConfig
		if !errors.As(err, &inErr) {
			t.Fatalf("expected error of type *ErrInClusterConfig via errors.As, got %T: %v", err, err)
		}
		if _, ok := err.(*ErrInClusterConfig); !ok {
			t.Fatalf("expected error to be *ErrInClusterConfig via type assertion, got %T", err)
		}

		// Unwrapped error should come from rest.InClusterConfig
		orig := errors.Unwrap(err)
		if orig == nil {
			t.Error("expected ErrInClusterConfig to wrap an underlying error, but Unwrap() returned nil")
		}
	})

	t.Run("In cluster success", func(t *testing.T) {
		// Simulate “inside cluster” by setting the marker env var.
		t.Setenv(InClusterEnv, "true")

		defer resetStubs()
		var inClusterConfigFuncCalled bool
		inClusterConfigFunc = func() (*rest.Config, error) {
			inClusterConfigFuncCalled = true
			return &rest.Config{}, nil
		}

		cfg, err := getClusterConfig()
		if cfg == nil {
			t.Fatal("expected non-nil config on in-cluster success")
		}
		if err != nil {
			t.Fatalf("expected nil error for inClusterConfigFunc, got: %v", err)
		}

		if !inClusterConfigFuncCalled {
			t.Fatal("expected inClusterConfigFunc to be called")
		}
	})

	t.Run("Out of cluster missing env", func(t *testing.T) {
		// Simulate “outside cluster” and missing HOME/USERPROFILE
		t.Setenv(InClusterEnv, "")
		t.Setenv(UnixEnvHome, "")
		t.Setenv(WindowsEnvHome, "")

		cfg, err := getClusterConfig()
		if cfg != nil {
			t.Fatalf("expected nil config when home env missing, got %+v", cfg)
		}
		if err == nil {
			t.Fatal("expected an error when out-of-cluster missing home env")
		}

		// Should propagate *ErrMissingEnvVars
		var missErr *ErrMissingEnvVars
		if !errors.As(err, &missErr) {
			t.Fatalf("expected error of type *ErrMissingEnvVars via errors.As, got %T: %v", err, err)
		}
		if _, ok := err.(*ErrMissingEnvVars); !ok {
			t.Fatalf("expected error to be *ErrMissingEnvVars via type assertion, got %T", err)
		}
	})

	t.Run("Out of cluster build failure", func(t *testing.T) {
		// Simulate “outside cluster” and a build failure in getOutOfClusterConfig
		//
		// We set HOME but do not create a kubeconfig file, so
		// clientcmd.BuildConfigFromFlags will error and getOut wraps it in *ErrOutClusterConfig.
		t.Setenv(InClusterEnv, "")
		t.Setenv(UnixEnvHome, UnixMockHome)
		t.Setenv(WindowsEnvHome, "")

		cfg, err := getClusterConfig()
		if cfg != nil {
			t.Fatalf("expected nil config on build failure, got %+v", cfg)
		}
		if err == nil {
			t.Fatal("expected an error when buildConfigFromFlags fails")
		}

		// Should be *ErrOutClusterConfig
		var outErr *ErrOutClusterConfig
		if !errors.As(err, &outErr) {
			t.Fatalf("expected error of type *ErrOutClusterConfig via errors.As, got %T: %v", err, err)
		}
		if _, ok := err.(*ErrOutClusterConfig); !ok {
			t.Fatalf("expected error to be *ErrOutClusterConfig via type assertion, got %T", err)
		}

		// Unwrap to the underlying error from BuildConfigFromFlags
		if errors.Unwrap(err) == nil {
			t.Error("expected ErrOutClusterConfig to wrap an underlying error, but Unwrap() returned nil")
		}
	})

	t.Run("Out of cluster success", func(t *testing.T) {
		// Simulate “outside cluster”
		t.Setenv(InClusterEnv, "")

		defer resetStubs()
		var outOfClusterConfigFuncCalled bool
		outOfClusterConfigFunc = func(_ ConfigProvider) (*rest.Config, error) {
			outOfClusterConfigFuncCalled = true
			return &rest.Config{}, nil
		}

		cfg, err := getClusterConfig()
		if cfg == nil {
			t.Fatal("expected non-nil config on out-cluster success")
		}
		if err != nil {
			t.Fatalf("expected nil error for outOfClusterConfigFunc, got: %v", err)
		}

		if !outOfClusterConfigFuncCalled {
			t.Fatal("expected outOfClusterConfigFunc to be called")
		}
	})
}

func TestCreateKubernetesClient(t *testing.T) {
	defer resetStubs()

	t.Run("Propegates known errors", func(t *testing.T) {
		tests := []struct {
			name       string
			clusterErr error // what getClusterConfigFunc should return
			expectType any   // pointer to the expected error type
		}{
			{
				name:       "MissingEnvVars",
				clusterErr: &ErrMissingEnvVars{},
				expectType: (*ErrMissingEnvVars)(nil),
			},
			{
				name:       "OutClusterConfig",
				clusterErr: &ErrOutClusterConfig{Err: errors.New("boom")},
				expectType: (*ErrOutClusterConfig)(nil),
			},
			{
				name:       "InClusterConfig",
				clusterErr: &ErrInClusterConfig{Err: errors.New("oops")},
				expectType: (*ErrInClusterConfig)(nil),
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Stub getClusterConfigFunc to return the specific error
				getClusterConfigFunc = func() (*rest.Config, error) {
					return nil, tc.clusterErr
				}

				client, err := CreateKubernetesClient()
				if client != nil {
					t.Errorf("expected nil client, got %+v", client)
				}
				if err == nil {
					t.Fatalf("expected error of type %T, got nil", tc.expectType)
				}

				// 1) Check via errors.As
				if !errors.As(err, &tc.expectType) {
					t.Fatalf("expected error to be %T via errors.As, got %T: %v", tc.expectType, err, err)
				}

				// 2) Direct type assertion
				switch tc.expectType.(type) {
				case *ErrMissingEnvVars:
					if _, ok := err.(*ErrMissingEnvVars); !ok {
						t.Fatalf("expected *ErrMissingEnvVars, got %T", err)
					}
				case *ErrOutClusterConfig:
					if _, ok := err.(*ErrOutClusterConfig); !ok {
						t.Fatalf("expected *ErrOutClusterConfig, got %T", err)
					}
				case *ErrInClusterConfig:
					if _, ok := err.(*ErrInClusterConfig); !ok {
						t.Fatalf("expected *ErrInClusterConfig, got %T", err)
					}
				}
			})
		}
	})

	t.Run("Panics on enexpected errors", func(t *testing.T) {
		// Stub getClusterConfigFunc to return an unexpected error
		getClusterConfigFunc = func() (*rest.Config, error) {
			return nil, fmt.Errorf("something else")
		}

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on unexpected error, but did not panic")
			}
		}()
		_, _ = CreateKubernetesClient()
	})

	t.Run("Client creation error", func(t *testing.T) {
		// 1) Stub getClusterConfigFunc to succeed
		fakeCfg := &rest.Config{Host: "dummy"}
		getClusterConfigFunc = func() (*rest.Config, error) {
			return fakeCfg, nil
		}

		// 2) Stub newForConfigFunc to return an error
		inner := errors.New("dial failed")
		newForConfigFunc = func(cfg *rest.Config) (*kubernetes.Clientset, error) {
			return nil, inner
		}

		client, err := CreateKubernetesClient()
		if client != nil {
			t.Errorf("expected nil client on creation error, got %+v", client)
		}
		if err == nil {
			t.Fatal("expected ErrClientCreation, got nil")
		}

		// 1) Check via errors.As
		var ccErr *ErrClientCreation
		if !errors.As(err, &ccErr) {
			t.Fatalf("expected error of type *ErrClientCreation via errors.As, got %T: %v", err, err)
		}
		// 2) Direct type assertion
		if _, ok := err.(*ErrClientCreation); !ok {
			t.Fatalf("expected *ErrClientCreation, got %T", err)
		}
		// 3) Unwrap
		if !errors.Is(err, inner) {
			t.Errorf("expected wrapped error to be inner; errors.Is(err, inner)=false")
		}
	})

	t.Run("Success", func(t *testing.T) {
		// 1) Stub getClusterConfigFunc to succeed
		fakeCfg := &rest.Config{Host: "dummy"}
		getClusterConfigFunc = func() (*rest.Config, error) {
			return fakeCfg, nil
		}

		// 2) Stub newForConfigFunc to return a fake clientset
		fakeClient := &kubernetes.Clientset{}
		newForConfigFunc = func(cfg *rest.Config) (*kubernetes.Clientset, error) {
			if cfg != fakeCfg {
				t.Fatalf("expected config %v, got %v", fakeCfg, cfg)
			}
			return fakeClient, nil
		}

		client, err := CreateKubernetesClient()
		if err != nil {
			t.Fatalf("unexpected error in happy path: %v", err)
		}
		if client != fakeClient {
			t.Fatalf("expected client %v, got %v", fakeClient, client)
		}
	})
}
