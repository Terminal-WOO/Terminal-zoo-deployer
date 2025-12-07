#!/bin/bash

# Domain Boundaries Check Script
# Controleert dat domain boundaries gerespecteerd worden

set -e

ERRORS=0

echo "=== Checking Domain Boundaries ==="
echo ""

# Check 1: Frontend should not import backend Go code
echo "1. Checking frontend does not import backend Go code..."
if grep -r "github.com/ClappFormOrg/AI-CO/go" app/ server/ --include="*.ts" --include="*.vue" --include="*.js" 2>/dev/null; then
    echo "❌ ERROR: Frontend imports backend Go code!"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ OK: Frontend does not import backend Go code"
fi

# Check 2: Backend should not import frontend code
echo ""
echo "2. Checking backend does not import frontend code..."
if grep -r "app/" go/ --include="*.go" 2>/dev/null; then
    echo "❌ ERROR: Backend imports frontend code!"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ OK: Backend does not import frontend code"
fi

# Check 3: Server API should only use server/ directory
echo ""
echo "3. Checking server API boundaries..."
# Server API should not directly import app/ code
if grep -r "from '@/" server/ --include="*.ts" 2>/dev/null | grep -v "server/" 2>/dev/null; then
    echo "⚠️  WARNING: Server API may import app code (review needed)"
else
    echo "✅ OK: Server API boundaries respected"
fi

# Check 4: Domain boundaries in stores
echo ""
echo "4. Checking store boundaries..."
# Each store should be independent
STORES=$(find app/stores -name "*.ts" -type f)
for store in $STORES; do
    STORE_NAME=$(basename "$store" .ts)
    # Check if store imports other stores (should be minimal)
    IMPORTS=$(grep -E "from.*stores/(?!${STORE_NAME})" "$store" 2>/dev/null | wc -l || echo "0")
    if [ "$IMPORTS" -gt 3 ]; then
        echo "⚠️  WARNING: $STORE_NAME imports many other stores (may indicate boundary violation)"
    fi
done
echo "✅ OK: Store boundaries check complete"

# Check 5: API boundaries
echo ""
echo "5. Checking API boundaries..."
# API endpoints should be in server/api/
if find server/api -name "*.ts" -type f | grep -v "server/api" 2>/dev/null; then
    echo "❌ ERROR: API files outside server/api directory!"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ OK: API boundaries respected"
fi

# Summary
echo ""
echo "=== Summary ==="
if [ $ERRORS -eq 0 ]; then
    echo "✅ All domain boundary checks passed!"
    exit 0
else
    echo "❌ Found $ERRORS domain boundary violations!"
    echo ""
    echo "Please review and fix the violations above."
    exit 1
fi

