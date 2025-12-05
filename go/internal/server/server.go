package server

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/ClappFormOrg/AI-CO/go/pkg/kube/client"
	"github.com/ClappFormOrg/AI-CO/go/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Component string = "internal.server"
)

type DomainConfig struct {
	Domain      string `json:"domain"`
	Certificate []byte `json:"certificate"` // PEM encoded
	PrivateKey  []byte `json:"privateKey"`  // PEM encoded
}

type Handler struct {
	mux            *http.ServeMux
	logger         log.Logger
	clientset      *kubernetes.Clientset
	clientctrl     ctrlclient.Client
	clients        map[string]*kubernetes.Clientset
	clientsConfig  map[string]*rest.Config
	clientsDomains map[string]DomainConfig
	tlsKey         []byte // WARN: Check for emptiness before use!
	tlsCrt         []byte // WARN: Check for emptiness before use!
}

func int32Ptr(i int32) *int32 { return &i }

func validateNamespaceExists(clientset *kubernetes.Clientset, namespace string) error {
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("namespace %s does not exist: %w", namespace, err)
	}
	return nil
}

// AuthMiddleware checks for a valid token in the Authorization header
func AuthMiddleware(next http.HandlerFunc, validToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || token != validToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func switchClientset(h *Handler, clusterName string) (*kubernetes.Clientset, error) {
	clientset, exists := h.clients[clusterName]
	if !exists {
		return nil, fmt.Errorf("unknown cluster: %s", clusterName)
	}
	return clientset, nil
}

type DeploymentRequest struct {
	Namespace      string `json:"namespace"`
	DeploymentName string `json:"deploymentName"`
	Image          string `json:"image"`
	Replicas       int32  `json:"replicas"`
	Resources      struct {
		CPULimits      string `json:"cpuLimits"`
		CPURequests    string `json:"cpuRequests"`
		MemoryLimits   string `json:"memoryLimits"`
		MemoryRequests string `json:"memoryRequests"`
	} `json:"resources"`
	Ports []corev1.ContainerPort `json:"ports"`
}

// validateDeploymentRequestBody checks all required fields in the deployment request body
func validateDeploymentRequestBody(req DeploymentRequest) error {
	var err []error
	if req.Namespace == "" {
		err = append(err, fmt.Errorf("namespace is required"))
	}
	if req.DeploymentName == "" {
		err = append(err, fmt.Errorf("deploymentName is required"))
	}
	if req.Image == "" {
		err = append(err, fmt.Errorf("image is required"))
	}
	if req.Replicas <= 0 {
		err = append(err, fmt.Errorf("replicas must be greater than 0"))
	}
	if req.Resources.CPULimits == "" {
		err = append(err, fmt.Errorf("resources.cpuLimits is required"))
	}
	if req.Resources.CPURequests == "" {
		err = append(err, fmt.Errorf("resources.cpuRequests is required"))
	}
	if req.Resources.MemoryLimits == "" {
		err = append(err, fmt.Errorf("resources.memoryLimits is required"))
	}
	if req.Resources.MemoryRequests == "" {
		err = append(err, fmt.Errorf("resources.memoryRequests is required"))
	}
	if len(req.Ports) == 0 {
		err = append(err, fmt.Errorf("at least one port is required in ports"))
	}
	for i, port := range req.Ports {
		if port.ContainerPort == 0 {
			err = append(err, fmt.Errorf("ports[%d].containerPort is required and must be > 0", i))
		}
	}
	if len(err) > 0 {
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}

func NewHandler(opts ...Option) (h *Handler, err error) {
	h = &Handler{
		mux:            new(http.ServeMux),
		logger:         log.NewNoOpLogger(),
		clients:        make(map[string]*kubernetes.Clientset),
		clientsConfig:  make(map[string]*rest.Config),
		clientsDomains: make(map[string]DomainConfig),
		tlsKey:         []byte{},
		tlsCrt:         []byte{},
	}

	for _, opt := range opts {
		opt(h)
	}

	// Create the main clientset
	clientset, clientConfig, err := client.CreateKubernetesClient()
	if err != nil {
		message := "failed to create kubernetes clientset"
		h.logger.Error(message, "err", err)
		return nil, fmt.Errorf("%s: %w", message, err)
	}
	// Set as active clientset
	h.clientset = clientset

	// Add to selectable clients list and the config
	h.clients["clappform"] = clientset
	h.clientsConfig["clappform"] = clientConfig

	h.clientctrl, err = client.CreateControllerRuntimeClient()
	if err != nil {
		message := "failed to create kubernetes controller runtime client"
		h.logger.Error(message, "err", err)
		return nil, fmt.Errorf("%s: %w", message, err)
	}

	print("Handler initialized successfully, Following clusters are configured:\n")
	for clusterName := range h.clients {
		print(fmt.Sprintf(" - %s\n", clusterName))
	}

	h.mux.HandleFunc("GET /pods/{namespace}", AuthMiddleware(h.handleActivePods(), ""))
	h.mux.HandleFunc("GET /pods/{namespace}/{podname}/logs", AuthMiddleware(h.handlePodLogs(), ""))

	h.mux.HandleFunc("GET /deployments/{namespace}", AuthMiddleware(h.handleDeploymentGetAll(), ""))
	h.mux.HandleFunc("GET /deployments/{namespace}/{deploymentName}", AuthMiddleware(h.handleDeploymentGet(), ""))
	h.mux.HandleFunc("POST /deployments", AuthMiddleware(h.handleDeploymentCreation(), ""))
	h.mux.HandleFunc("DELETE /deployments/{namespace}/{deploymentName}", AuthMiddleware(h.handleDeploymentDeletion(), ""))
	h.mux.HandleFunc("PUT /deployments/{namespace}/{deploymentName}", AuthMiddleware(h.handleDeploymentUpdate(), ""))
	h.mux.HandleFunc("POST /deployments/{namespace}/{deploymentName}/restart", AuthMiddleware(h.handleRolloutRestart(), ""))

	h.mux.HandleFunc("GET /clusters", AuthMiddleware(h.handleListClusters(), ""))
	h.mux.HandleFunc("GET /clusters/{clusterName}", AuthMiddleware(h.handleGetCluster(), ""))
	h.mux.HandleFunc("POST /clusters", AuthMiddleware(h.handleAddClusterContext(), ""))

	h.mux.HandleFunc("POST /secrets", AuthMiddleware(h.handleCreateSecret(), ""))
	h.mux.HandleFunc("GET /secrets/{namespace}", AuthMiddleware(h.handleGetSecrets(), ""))

	h.mux.HandleFunc("POST /configmap", AuthMiddleware(h.handleCreateConfigMap(), ""))
	h.mux.HandleFunc("GET /configmap/{namespace}", AuthMiddleware(h.handleGetConfigMaps(), ""))

	// Health check endpoints (no auth required)
	h.mux.HandleFunc("GET /health", h.handleHealth())
	h.mux.HandleFunc("GET /ready", h.handleReady())

	return h, nil
}

func (h *Handler) Stop(ctx context.Context) error {
	h.logger.DebugCtx(ctx, "stopping handler")
	return nil
}

func (h *Handler) Close() error {
	h.logger.Debug("closing handler")
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) handleActivePods() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace from path
		namespace := r.PathValue("namespace")
		if namespace == "" {
			http.Error(w, "namespace is empty", http.StatusBadRequest)
			return
		}

		// Determine which clientset to use
		activeClientset := h.clientset
		clusterName := r.Header.Get("cluster-name")
		if clusterName != "" {
			var err error
			activeClientset, err = switchClientset(h, clusterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Implementation for handling active pods
		pods, err := activeClientset.CoreV1().Pods(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list pods: %v", err), http.StatusInternalServerError)
			return
		}

		// return all information about the pods in json format
		podInfos := make([]map[string]interface{}, 0, len(pods.Items))
		for _, pod := range pods.Items {
			podInfo := map[string]interface{}{
				"name":             pod.Name,
				"status":           pod.Status.Phase,
				"startTime":        pod.Status.StartTime,
				"restartPolicy":    pod.Spec.RestartPolicy,
				"statusConditions": pod.Status.Conditions,
				"resources":        pod.Spec.Containers,
			}
			podInfos = append(podInfos, podInfo)
		}

		// Set response header and write JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pods": podInfos,
		})
	}
}

func (h *Handler) handlePodLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace and podName from path
		namespace := r.PathValue("namespace")
		podName := r.PathValue("podname")

		if namespace == "" || podName == "" {
			http.Error(w, "namespace and podName are required", http.StatusBadRequest)
			return
		}

		// Get the logs for the specified pod
		podLogOpts := corev1.PodLogOptions{}
		req := h.clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
		podLogs, err := req.Stream(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get pod logs: %v", err), http.StatusInternalServerError)
			return
		}
		defer podLogs.Close()
		buf := make([]byte, 2000)
		numBytes, err := podLogs.Read(buf)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read pod logs: %v", err), http.StatusInternalServerError)
			return
		}
		logs := string(buf[:numBytes])
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(logs))
	}
}

func (h *Handler) handleRolloutRestart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace and deploymentName from path
		namespace := r.PathValue("namespace")
		if namespace == "" {
			http.Error(w, "namespace is empty", http.StatusBadRequest)
			return
		}
		deploymentName := r.PathValue("deploymentName")
		if deploymentName == "" {
			http.Error(w, "deploymentName is empty", http.StatusBadRequest)
			return
		}

		// Implementation for handling rollout restart
		deploymentsClient := h.clientset.AppsV1().Deployments(namespace)
		deployment, err := deploymentsClient.Get(r.Context(), deploymentName, metav1.GetOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get deployment: %v", err), http.StatusInternalServerError)
			return
		}

		// Patch the deployment with a new annotation to trigger a rollout restart
		if deployment.Spec.Template.Annotations == nil {
			deployment.Spec.Template.Annotations = make(map[string]string)
		}
		deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().UTC().Format(time.RFC3339)
		_, err = deploymentsClient.Update(r.Context(), deployment, metav1.UpdateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update deployment: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("deployment restarted successfully"))
	}
}

func (h *Handler) handleDeploymentCreation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 0) decode + validate
		var in DeploymentRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode body: %v", err), http.StatusBadRequest)
			return
		}
		if err := validateDeploymentRequestBody(in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(in.Ports) == 0 {
			http.Error(w, "at least one container port is required", http.StatusBadRequest)
			return
		}

		// 1) pick cluster client
		cs := h.clientset
		domainConfig := DomainConfig{
			Domain:      "services.clappform.com",
			Certificate: h.tlsCrt,
			PrivateKey:  h.tlsKey,
		}
		clientConfig := h.clientsConfig["clappform"]
		if name := r.Header.Get("cluster-name"); name != "" {
			var err error
			cs, err = switchClientset(h, name)
			clientConfig = h.clientsConfig[name]
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if dc, ok := h.clientsDomains[name]; ok {
				domainConfig = dc
			}
		}

		// 2) ensure namespace exists
		if err := validateNamespaceExists(cs, in.Namespace); err != nil {
			if _, err := cs.CoreV1().Namespaces().Create(r.Context(), &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: in.Namespace},
			}, metav1.CreateOptions{}); err != nil {
				http.Error(w, fmt.Sprintf("failed to create namespace: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// 3) ensure TLS secret (only if you really need TLS now)
		// NOTE: replace these with real cert/key data or skip TLS until ready.
		if _, err := cs.CoreV1().Secrets(in.Namespace).Get(r.Context(), "cert", metav1.GetOptions{}); err != nil {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "cert", Namespace: in.Namespace},
				Type:       corev1.SecretTypeTLS,
				Data: map[string][]byte{
					"tls.crt": slices.Clone(domainConfig.Certificate),
					"tls.key": slices.Clone(domainConfig.PrivateKey),
				},
			}
			if _, err := cs.CoreV1().Secrets(in.Namespace).Create(r.Context(), secret, metav1.CreateOptions{}); err != nil {
				http.Error(w, fmt.Sprintf("failed to create TLS secret: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Common names & labels
		depName := in.DeploymentName + "-deployment"
		svcName := in.DeploymentName + "-deployment-service"
		appLabel := map[string]string{"app": in.DeploymentName}

		// 4) Deployment
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      depName,
				Namespace: in.Namespace,
				Labels:    appLabel,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(in.Replicas),
				Selector: &metav1.LabelSelector{MatchLabels: appLabel},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: appLabel},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  in.DeploymentName + "-container",
								Image: in.Image,
								Ports: in.Ports, // []corev1.ContainerPort
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(in.Resources.CPULimits),
										corev1.ResourceMemory: resource.MustParse(in.Resources.MemoryLimits),
									},
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse(in.Resources.CPURequests),
										corev1.ResourceMemory: resource.MustParse(in.Resources.MemoryRequests),
									},
								},
								// Optional but recommended:
								// ReadinessProbe: ...,
								// LivenessProbe:  ...,
							},
						},
					},
				},
			},
		}

		depClient := cs.AppsV1().Deployments(in.Namespace)
		createdDep, err := depClient.Create(r.Context(), deployment, metav1.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create deployment: %v", err), http.StatusInternalServerError)
			return
		}

		// 5) Service (name must match Ingress backend!)
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      svcName,
				Namespace: in.Namespace,
				Labels:    appLabel,
			},
			Spec: corev1.ServiceSpec{
				Selector: appLabel,
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Port:       in.Ports[0].ContainerPort,                      // service port
						TargetPort: intstr.FromInt(int(in.Ports[0].ContainerPort)), // forward to container
						// Name: optional but good if you add more ports
					},
				},
				Type: corev1.ServiceTypeClusterIP,
			},
		}
		if _, err := cs.CoreV1().Services(in.Namespace).Create(r.Context(), svc, metav1.CreateOptions{}); err != nil {
			http.Error(w, fmt.Sprintf("failed to create service: %v", err), http.StatusInternalServerError)
			return
		}

		// 6) Ingress (Traefik)
		// Use spec.IngressClassName and make sure traefik is installed & watching this namespace.
		// ingressClass := "traefik"
		pathPrefix := "/" + depName // if you want /<deployment>-deployment
		middleWareName := "strip-" + depName + "-prefix"

		gvr := schema.GroupVersionResource{Group: "traefik.io", Version: "v1alpha1", Resource: "middlewares"}
		dc, _ := dynamic.NewForConfig(clientConfig)

		obj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "traefik.io/v1alpha1",
				"kind":       "Middleware",
				"metadata":   map[string]interface{}{"name": middleWareName, "namespace": in.Namespace},
				"spec": map[string]interface{}{
					"stripPrefixRegex": map[string]interface{}{"regex": []string{"^/" + depName}},
				},
			},
		}

		response, err := dc.Resource(gvr).Namespace(in.Namespace).Create(r.Context(), obj, metav1.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create middleware: %v", err), http.StatusInternalServerError)
			return
		}

		_ = response // avoid unused var warning, though you might want to log or inspect it

		print("Domain config:\n")
		print(fmt.Sprintf(" - Domain: %s\n", domainConfig.Domain))
		print(fmt.Sprintf(" - Cert: %d bytes\n", len(domainConfig.Certificate)))
		print(fmt.Sprintf(" - Key: %d bytes\n", len(domainConfig.PrivateKey)))

		ing := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      in.DeploymentName + "-ingress",
				Namespace: in.Namespace,
				Annotations: map[string]string{
					"traefik.ingress.kubernetes.io/router.entrypoints": "websecure",
					"traefik.ingress.kubernetes.io/router.tls":         "true",
					"traefik.ingress.kubernetes.io/router.middlewares": fmt.Sprintf("%s-%s@kubernetescrd", in.Namespace, middleWareName),
				},
			},
			Spec: networkingv1.IngressSpec{
				// IngressClassName: &ingressClass,
				Rules: []networkingv1.IngressRule{
					{
						Host: domainConfig.Domain,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     pathPrefix,
										PathType: ptrPathType(networkingv1.PathTypePrefix),
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: svcName, // <-- must match the Service name
												Port: networkingv1.ServiceBackendPort{
													Number: in.Ports[0].ContainerPort,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				TLS: []networkingv1.IngressTLS{
					{Hosts: []string{"services.clappform.com"}, SecretName: "cert"},
				},
			},
		}
		if _, err := cs.NetworkingV1().Ingresses(in.Namespace).Create(r.Context(), ing, metav1.CreateOptions{}); err != nil {
			http.Error(w, fmt.Sprintf("failed to create ingress: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(createdDep)
	}
}

func ptrPathType(pt networkingv1.PathType) *networkingv1.PathType { return &pt }

func (h *Handler) handleDeploymentDeletion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace and deploymentName from path
		namespace := r.PathValue("namespace")
		deploymentName := r.PathValue("deploymentName")
		if namespace == "" || deploymentName == "" {
			http.Error(w, "namespace and deploymentName are required", http.StatusBadRequest)
			return
		}

		// Determine which clientset to use
		activeClientset := h.clientset
		clusterName := r.Header.Get("cluster-name")
		if clusterName != "" {
			var err error
			activeClientset, err = switchClientset(h, clusterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Validate namespace exists
		if err := validateNamespaceExists(activeClientset, namespace); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Delete the specified deployment
		err := activeClientset.AppsV1().Deployments(namespace).Delete(r.Context(), deploymentName, metav1.DeleteOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to delete deployment: %v", err), http.StatusInternalServerError)
			return
		}

		// Also delete the associated service
		err = activeClientset.CoreV1().Services(namespace).Delete(r.Context(), deploymentName+"-service", metav1.DeleteOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to delete service: %v", err), http.StatusInternalServerError)
			return
		}

		// Also delete the associated ingress
		err = activeClientset.NetworkingV1().Ingresses(namespace).Delete(r.Context(), deploymentName+"-ingress", metav1.DeleteOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to delete ingress: %v", err), http.StatusInternalServerError)
			return
		}

		// Check if there are any deployments left in the namespace otherwise delete the namespace
		deployments, err := activeClientset.AppsV1().Deployments(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list deployments: %v", err), http.StatusInternalServerError)
			return
		}
		if len(deployments.Items) == 0 {
			err = activeClientset.CoreV1().Namespaces().Delete(r.Context(), namespace, metav1.DeleteOptions{})
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to delete namespace: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Return success response
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) handleDeploymentUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace and deploymentName from path
		namespace := r.PathValue("namespace")
		deploymentName := r.PathValue("deploymentName")
		if namespace == "" || deploymentName == "" {
			http.Error(w, "namespace and deploymentName are required", http.StatusBadRequest)
			return
		}

		// Load data from request body
		var requestBody struct {
			Replicas int32 `json:"replicas"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}
		if requestBody.Replicas <= 0 {
			http.Error(w, "replicas must be greater than 0", http.StatusBadRequest)
			return
		}
		// Get the existing deployment
		deploymentsClient := h.clientset.AppsV1().Deployments(namespace)
		deployment, err := deploymentsClient.Get(r.Context(), deploymentName, metav1.GetOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get deployment: %v", err), http.StatusInternalServerError)
			return
		}
		// Update the replicas
		deployment.Spec.Replicas = int32Ptr(requestBody.Replicas)
		updatedDeployment, err := deploymentsClient.Update(r.Context(), deployment, metav1.UpdateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update deployment: %v", err), http.StatusInternalServerError)
			return
		}
		// Return the updated deployment in JSON format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedDeployment)
	}
}

func (h *Handler) handleDeploymentGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace and deploymentName from path
		namespace := r.PathValue("namespace")
		deploymentName := r.PathValue("deploymentName")
		if namespace == "" || deploymentName == "" {
			http.Error(w, "namespace and deploymentName are required", http.StatusBadRequest)
			return
		}

		// Determine which clientset to use
		activeClientset := h.clientset
		clusterName := r.Header.Get("cluster-name")
		if clusterName != "" {
			var err error
			activeClientset, err = switchClientset(h, clusterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Implementation for handling deployment retrieval
		deployment, err := activeClientset.AppsV1().Deployments(namespace).Get(r.Context(), deploymentName, metav1.GetOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get deployment: %v", err), http.StatusInternalServerError)
			return
		}

		// Return the deployment in JSON format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deployment)
	}
}

func (h *Handler) handleDeploymentGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace from path
		namespace := r.PathValue("namespace")
		if namespace == "" {
			http.Error(w, "namespace is required", http.StatusBadRequest)
			return
		}

		// Determine which clientset to use
		activeClientset := h.clientset
		clusterName := r.Header.Get("cluster-name")
		if clusterName != "" {
			var err error
			activeClientset, err = switchClientset(h, clusterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Implementation for handling deployment retrieval
		deployments, err := activeClientset.AppsV1().Deployments(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list deployments: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deployments)
	}
}

// CLUSTERS
func (h *Handler) handleListClusters() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type clusterInfo struct {
			Name       string `json:"name"`
			Server     string `json:"server"`
			K8sVersion string `json:"k8sVersion"`
		}

		var clusters []clusterInfo
		for name, cfg := range h.clientsConfig {
			if name == "clappform" {
				continue // skip the main cluster
			}
			// Create a new clientset for the cluster
			cs, err := kubernetes.NewForConfig(cfg)
			if err != nil {
				continue // skip unreachable clusters
			}

			// get server version
			info, err := cs.Discovery().ServerVersion()
			if err != nil {
				http.Error(w, "failed to get server version", http.StatusInternalServerError)
				continue
			}

			out := clusterInfo{
				Name:       name,
				Server:     cfg.Host,
				K8sVersion: info.GitVersion,
			}

			clusters = append(clusters, out)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(clusters)
	}
}

func (h *Handler) handleGetCluster() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clusterName := r.PathValue("clusterName")
		if clusterName == "" {
			http.Error(w, "clusterName is required", http.StatusBadRequest)
			return
		}

		cfg, ok := h.clientsConfig[clusterName]
		if !ok {
			http.Error(w, "cluster not found", http.StatusNotFound)
			return
		}

		cs, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			http.Error(w, "failed to create client", http.StatusInternalServerError)
			return
		}

		info, err := cs.Discovery().ServerVersion()
		if err != nil {
			http.Error(w, "failed to get server version", http.StatusInternalServerError)
			return
		}

		out := struct {
			Name       string `json:"name"`
			Server     string `json:"server"`
			K8sVersion string `json:"k8sVersion"`
		}{
			Name:       clusterName,
			Server:     cfg.Host,
			K8sVersion: info.GitVersion,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}
}

func (h *Handler) handleCreateSecret() http.HandlerFunc {
	type req struct {
		Namespace  string            `json:"namespace"`
		SecretName string            `json:"secretName"`
		Data       map[string][]byte `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Get body
		var in req
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate secretName and data
		if in.SecretName == "" {
			http.Error(w, "secretName is required", http.StatusBadRequest)
			return
		}
		if len(in.Data) == 0 {
			http.Error(w, "data is required", http.StatusBadRequest)
			return
		}

		// Determine which clientset to use
		activeClientset := h.clientset
		clusterName := r.Header.Get("cluster-name")
		if clusterName != "" {
			var err error
			activeClientset, err = switchClientset(h, clusterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Validate namespace exists
		if err := validateNamespaceExists(activeClientset, in.Namespace); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Implementation for creating a secret
		secretName := in.SecretName
		if secretName == "" {
			http.Error(w, "secretName is required", http.StatusBadRequest)
			return
		}

		var secretData = make(map[string]string)
		for key, value := range in.Data {
			secretData[key] = string(value)
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: in.Namespace,
			},
			StringData: secretData,
		}
		_, err := activeClientset.CoreV1().Secrets(in.Namespace).Create(r.Context(), secret, metav1.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create secret: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) handleGetSecrets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace and secretName from path
		namespace := r.PathValue("namespace")
		if namespace == "" {
			http.Error(w, "namespace is required", http.StatusBadRequest)
			return
		}

		// Determine which clientset to use
		activeClientset := h.clientset
		clusterName := r.Header.Get("cluster-name")
		if clusterName != "" {
			var err error
			activeClientset, err = switchClientset(h, clusterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Implementation for getting secrets
		secrets, err := activeClientset.CoreV1().Secrets(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list secrets: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(secrets)

	}
}

func (h *Handler) handleCreateConfigMap() http.HandlerFunc {
	type req struct {
		Namespace     string            `json:"namespace"`
		ConfigMapName string            `json:"configMapName"`
		Data          map[string]string `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		// Get body
		var in req
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// Determine which clientset to use
		activeClientset := h.clientset
		clusterName := r.Header.Get("cluster-name")
		if clusterName != "" {
			var err error
			activeClientset, err = switchClientset(h, clusterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Validate namespace exists
		if err := validateNamespaceExists(activeClientset, in.Namespace); err != nil {
			// Create the namespace if it does not exist
			if _, err := activeClientset.CoreV1().Namespaces().Create(r.Context(), &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: in.Namespace},
			}, metav1.CreateOptions{}); err != nil {
				http.Error(w, fmt.Sprintf("failed to create namespace: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Implementation for creating a config map
		configMapName := in.ConfigMapName
		if configMapName == "" {
			http.Error(w, "configMapName is required", http.StatusBadRequest)
			return
		}

		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: in.Namespace,
			},
			Data: in.Data,
		}
		_, err := activeClientset.CoreV1().ConfigMaps(in.Namespace).Create(r.Context(), configMap, metav1.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create config map: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) handleGetConfigMaps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load namespace from path
		namespace := r.PathValue("namespace")
		if namespace == "" {
			http.Error(w, "namespace is required", http.StatusBadRequest)
			return
		}
		// Determine which clientset to use
		activeClientset := h.clientset
		clusterName := r.Header.Get("cluster-name")
		if clusterName != "" {
			var err error
			activeClientset, err = switchClientset(h, clusterName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Implementation for getting config maps
		configMaps, err := activeClientset.CoreV1().ConfigMaps(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list config maps: %v", err), http.StatusInternalServerError)
			return
		}

		// Remove the kube-root-ca.crt from the list if present
		var filteredConfigMaps []corev1.ConfigMap
		for _, cm := range configMaps.Items {
			if cm.Name == "kube-root-ca.crt" {
				continue
			}
			filteredConfigMaps = append(filteredConfigMaps, cm)
		}

		// Return the filtered list
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filteredConfigMaps)
	}
}

func (h *Handler) handleAddClusterContext() http.HandlerFunc {
	type req struct {
		Name             string `json:"name"`
		Server           string `json:"server"`
		CAPEM            string `json:"caPEM"`       // PEM string OR base64-encoded PEM
		BearerToken      string `json:"bearerToken"` // SA token
		DefaultNamespace string `json:"defaultNamespace,omitempty"`
		Domain           string `json:"domain,omitempty"`
		Certificate      []byte `json:"certificate,omitempty"`
		PrivateKey       []byte `json:"privateKey,omitempty"`
	}
	type resp struct {
		Name      string `json:"name"`
		Server    string `json:"server"`
		Version   string `json:"k8sVersion"`
		Namespace string `json:"defaultNamespace,omitempty"`
	}

	nameRe := regexp.MustCompile(`^[a-zA-Z0-9._-]{1,80}$`)

	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
		var in req
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		print("Received onboarding request with parameters:\n")
		fmt.Printf("%+v\n", in)

		print(fmt.Sprintf(" - Domain: %s\n", in.Domain))
		print(fmt.Sprintf(" - Cert: %d bytes\n", len(in.Certificate)))
		print(fmt.Sprintf(" - Key: %d bytes\n", len(in.PrivateKey)))

		if !nameRe.MatchString(in.Name) {
			http.Error(w, "invalid name", http.StatusBadRequest)
			return
		}
		if !strings.HasPrefix(in.Server, "https://") {
			http.Error(w, "server must start with https://", http.StatusBadRequest)
			return
		}
		if in.CAPEM == "" || in.BearerToken == "" {
			http.Error(w, "caPEM and bearerToken are required", http.StatusBadRequest)
			return
		}

		caBytes, err := normalizePEM(in.CAPEM)
		if err != nil {
			http.Error(w, "invalid caPEM", http.StatusBadRequest)
			return
		}

		// Check if we already have a client for this name
		if _, exists := h.clients[in.Name]; exists {
			http.Error(w, "cluster with this name already exists", http.StatusConflict)
			return
		}

		print("Onboarding new cluster:\n")
		print(fmt.Sprintf(" - Name: %s\n", in.Name))
		print(fmt.Sprintf(" - Server: %s\n", in.Server))
		print(fmt.Sprintf(" - CA: %d bytes\n", len(caBytes)))
		print(fmt.Sprintf(" - Token: %d bytes\n", len(in.BearerToken)))
		print(fmt.Sprintf(" - Default Namespace: %s\n", in.DefaultNamespace))

		// Build the Domain config if provided and store in map
		if (in.Domain != "" && (in.Certificate == nil || in.PrivateKey == nil)) ||
			(in.Domain == "" && (in.Certificate != nil || in.PrivateKey != nil)) {
			http.Error(w, "domain, certificate and privateKey must be all provided or all omitted", http.StatusBadRequest)
			return
		} else {
			domainConfig := DomainConfig{
				Domain:      in.Domain,
				Certificate: in.Certificate,
				PrivateKey:  in.PrivateKey,
			}

			h.clientsDomains[in.Name] = domainConfig
		}

		// Build the rest.Config
		cfg := &rest.Config{
			Host:        in.Server,
			BearerToken: in.BearerToken,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: caBytes,
			},
			UserAgent: "myapp/sa-onboarder",
			// Add the domain to the config if provided (used for ingress rules)

			// Optional client-side rate limits:
			// QPS: 5, Burst: 10,
		}

		cs, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			http.Error(w, "failed to build client", http.StatusInternalServerError)
			return
		}

		// Probe /version with a short timeout
		_, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		info, err := cs.Discovery().ServerVersion()
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to contact cluster: %v", err), http.StatusBadGateway)
			return
		}

		// Store the client
		h.clients[in.Name] = cs
		h.clientsConfig[in.Name] = cfg

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp{
			Name:      in.Name,
			Server:    in.Server,
			Version:   info.GitVersion,
			Namespace: in.DefaultNamespace,
		})
	}
}

func normalizePEM(s string) ([]byte, error) {
	// Accept raw PEM or base64-encoded PEM
	raw := strings.TrimSpace(s)
	// if it looks like base64 without PEM headers, try to decode
	if !strings.HasPrefix(raw, "-----BEGIN") {
		dec, err := base64.StdEncoding.DecodeString(raw)
		if err == nil {
			raw = string(dec)
		}
	}
	// Basic validation: parse at least one cert
	var block *pem.Block
	data := []byte(raw)
	var found bool
	for {
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			if _, err := x509.ParseCertificate(block.Bytes); err == nil {
				found = true
			}
		}
	}
	if !found {
		return nil, errors.New("no valid certificate found in PEM")
	}
	return []byte(raw), nil
}

// handleHealth returns a simple health check endpoint
func (h *Handler) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	}
}

// handleReady returns a readiness check endpoint
// Checks if the Kubernetes clientset is available
func (h *Handler) handleReady() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if we can connect to Kubernetes
		if h.clientset == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "not ready",
				"reason": "kubernetes client not initialized",
			})
			return
		}

		// Try to get server version as a connectivity check
		_, err := h.clientset.Discovery().ServerVersion()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "not ready",
				"reason": fmt.Sprintf("kubernetes connection failed: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ready",
		})
	}
}
