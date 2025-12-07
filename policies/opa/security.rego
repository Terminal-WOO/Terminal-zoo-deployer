package platform.security

import rego.v1

# Default deny
default allow = false

# Allow if all security checks pass
allow {
    input.kind == "Deployment"
    validate_image_security
    validate_secrets
    validate_network_policy
}

# Validate image security
validate_image_security {
    container := input.spec.template.spec.containers[_]
    
    # Image must be from allowed registry
    startswith(container.image, "rg.nl-ams.scw.cloud/")
    
    # Image must have tag (not latest in production)
    contains(container.image, ":")
    not endswith(container.image, ":latest") or input.metadata.namespace != "nl-appstore-registry"
}

# Validate secrets usage
validate_secrets {
    # Secrets must be from Kubernetes secrets, not hardcoded
    container := input.spec.template.spec.containers[_]
    env := container.env[_]
    
    # Environment variables from secrets are allowed
    env.valueFrom.secretKeyRef != null or env.value == null
}

# Validate network policy compliance
validate_network_policy {
    # All pods must have labels for network policy matching
    input.metadata.labels.app != ""
    input.metadata.labels.namespace != ""
}

# Error messages
errors[msg] {
    input.kind == "Deployment"
    not validate_image_security
    msg := "Image must be from allowed registry and have proper tag"
}

errors[msg] {
    input.kind == "Deployment"
    not validate_secrets
    msg := "Secrets must be referenced from Kubernetes secrets, not hardcoded"
}

errors[msg] {
    input.kind == "Deployment"
    not validate_network_policy
    msg := "Deployment must have proper labels for network policy matching"
}

