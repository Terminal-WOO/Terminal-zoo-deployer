#!/bin/bash

# Kubernetes Manifest Validation Test
# Validates all Kubernetes manifests in k8s/ directory

set -e

echo "=== Kubernetes Manifest Validation ==="
echo ""

ERRORS=0

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install kubectl."
    exit 1
fi

# Check if kubeval is available (optional but recommended)
if command -v kubeval &> /dev/null; then
    echo "✅ kubeval found"
    USE_KUBEVAL=true
else
    echo "⚠️  kubeval not found, using kubectl dry-run only"
    USE_KUBEVAL=false
fi

# Validate each manifest
MANIFESTS=(
    "k8s/namespace.yaml"
    "k8s/configmap.yaml"
    "k8s/frontend-deployment.yaml"
    "k8s/backend-deployment.yaml"
    "k8s/frontend-service.yaml"
    "k8s/backend-service.yaml"
    "k8s/ingress.yaml"
    "k8s/cluster-issuer.yaml"
)

for manifest in "${MANIFESTS[@]}"; do
    if [ ! -f "$manifest" ]; then
        echo "⚠️  Skipping $manifest (not found)"
        continue
    fi

    echo "Validating $manifest..."
    
    # Validate with kubectl dry-run
    if kubectl apply --dry-run=client -f "$manifest" > /dev/null 2>&1; then
        echo "  ✅ kubectl validation passed"
    else
        echo "  ❌ kubectl validation failed"
        kubectl apply --dry-run=client -f "$manifest"
        ERRORS=$((ERRORS + 1))
    fi

    # Validate with kubeval if available
    if [ "$USE_KUBEVAL" = true ]; then
        if kubeval "$manifest" --strict > /dev/null 2>&1; then
            echo "  ✅ kubeval validation passed"
        else
            echo "  ❌ kubeval validation failed"
            kubeval "$manifest" --strict
            ERRORS=$((ERRORS + 1))
        fi
    fi
done

# Validate network policies
if [ -d "k8s/network-policies" ]; then
    echo ""
    echo "Validating network policies..."
    for policy in k8s/network-policies/*.yaml; do
        if [ -f "$policy" ]; then
            if kubectl apply --dry-run=client -f "$policy" > /dev/null 2>&1; then
                echo "  ✅ $(basename $policy) validation passed"
            else
                echo "  ❌ $(basename $policy) validation failed"
                kubectl apply --dry-run=client -f "$policy"
                ERRORS=$((ERRORS + 1))
            fi
        fi
    done
fi

# Summary
echo ""
echo "=== Summary ==="
if [ $ERRORS -eq 0 ]; then
    echo "✅ All manifests validated successfully!"
    exit 0
else
    echo "❌ Found $ERRORS validation errors!"
    exit 1
fi

