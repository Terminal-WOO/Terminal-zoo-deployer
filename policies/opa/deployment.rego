package platform.deployment

import rego.v1

# Default deny
default allow = false

# Allow deployment if all checks pass
allow {
    input.kind == "Deployment"
    validate_namespace
    validate_resources
    validate_security
    validate_labels
}

# Validate namespace exists and is allowed
validate_namespace {
    allowed_namespaces := {"nl-appstore-registry", "default"}
    allowed_namespaces[input.metadata.namespace]
}

# Validate resource limits are within allowed ranges
validate_resources {
    container := input.spec.template.spec.containers[_]
    
    # CPU limits
    container.resources.limits.cpu != ""
    parse_cpu(container.resources.limits.cpu) <= 4  # Max 4 CPU
    
    # Memory limits
    container.resources.limits.memory != ""
    parse_memory(container.resources.limits.memory) <= 8 * 1024 * 1024 * 1024  # Max 8GB
    
    # CPU requests
    container.resources.requests.cpu != ""
    parse_cpu(container.resources.requests.cpu) <= 2  # Max 2 CPU request
    
    # Memory requests
    container.resources.requests.memory != ""
    parse_memory(container.resources.requests.memory) <= 4 * 1024 * 1024 * 1024  # Max 4GB request
}

# Validate security constraints
validate_security {
    container := input.spec.template.spec.containers[_]
    
    # Must run as non-root
    input.spec.template.spec.securityContext.runAsNonRoot == true
    
    # Read-only root filesystem (if specified)
    not container.securityContext.readOnlyRootFilesystem or container.securityContext.readOnlyRootFilesystem == true
    
    # No privileged containers
    not container.securityContext.privileged or container.securityContext.privileged == false
}

# Validate required labels
validate_labels {
    input.metadata.labels.app != ""
    input.metadata.labels.version != ""
    input.metadata.labels.managed-by == "platform"
}

# Parse CPU string (e.g., "500m", "1", "2")
parse_cpu(cpu_str) = result {
    # Handle millicores (e.g., "500m" = 0.5)
    re_match("^[0-9]+m$", cpu_str)
    result := to_number(trim(cpu_str, "m")) / 1000
}

parse_cpu(cpu_str) = result {
    # Handle cores (e.g., "1", "2")
    re_match("^[0-9]+$", cpu_str)
    result := to_number(cpu_str)
}

# Parse memory string (e.g., "512Mi", "1Gi", "2G")
parse_memory(mem_str) = result {
    # Handle Mi (Mebibytes)
    re_match("^[0-9]+Mi$", mem_str)
    result := to_number(trim(mem_str, "Mi")) * 1024 * 1024
}

parse_memory(mem_str) = result {
    # Handle Gi (Gibibytes)
    re_match("^[0-9]+Gi$", mem_str)
    result := to_number(trim(mem_str, "Gi")) * 1024 * 1024 * 1024
}

parse_memory(mem_str) = result {
    # Handle G (Gigabytes)
    re_match("^[0-9]+G$", mem_str)
    result := to_number(trim(mem_str, "G")) * 1000 * 1000 * 1000
}

# Error messages
errors[msg] {
    input.kind == "Deployment"
    not validate_namespace
    msg := "Deployment namespace not allowed"
}

errors[msg] {
    input.kind == "Deployment"
    not validate_resources
    msg := "Resource limits exceed allowed maximums (CPU: 4, Memory: 8Gi)"
}

errors[msg] {
    input.kind == "Deployment"
    not validate_security
    msg := "Security constraints not met (must run as non-root)"
}

errors[msg] {
    input.kind == "Deployment"
    not validate_labels
    msg := "Required labels missing (app, version, managed-by)"
}

