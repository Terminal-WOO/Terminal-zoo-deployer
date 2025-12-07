# Software Supply Chain Security

Dit directory bevat policies en tools voor software supply chain security.

## Overzicht

Software supply chain security zorgt ervoor dat alle software components vertrouwd zijn en geen security vulnerabilities bevatten.

## Componenten

### 1. SBOM (Software Bill of Materials)

**Doel**: Automatische generatie van SBOM voor alle container images

**Tools**:
- **Syft**: SBOM generatie tool
- **SPDX/CycloneDX**: SBOM formaten

**Implementation**:
```bash
# Generate SBOM for frontend image
syft docker:rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest -o spdx-json > frontend-sbom.json

# Generate SBOM for backend image
syft docker:rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest -o spdx-json > backend-sbom.json
```

**Location**: `policies/supply-chain/sbom/`

---

### 2. Image Signing (Cosign)

**Doel**: Sign en verify container images

**Tools**:
- **Cosign**: Image signing tool van Sigstore

**Implementation**:
```bash
# Sign image
cosign sign --key cosign.key rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest

# Verify image signature
cosign verify --key cosign.pub rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest
```

**Location**: `policies/supply-chain/signing/`

---

### 3. Vulnerability Scanning

**Doel**: Automatische scanning van images en dependencies

**Tools**:
- **Trivy**: Vulnerability scanner
- **Docker Scout**: Docker's vulnerability scanner
- **npm audit**: Node.js dependency scanning
- **go list -u**: Go dependency checking

**Implementation**:
```bash
# Scan container image
trivy image rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest

# Scan dependencies
npm audit
go list -u -m all
```

**Location**: `policies/supply-chain/scanning/`

---

### 4. Dependency Tracking

**Doel**: Tracking van alle dependencies en updates

**Tools**:
- **Dependabot**: Automated dependency updates
- **Renovate**: Alternative dependency updater

**Configuration**:
- `.github/dependabot.yml` - Dependabot configuratie
- `renovate.json` - Renovate configuratie (alternatief)

---

## Supply Chain Security Workflow

### Build Time

1. **Dependency Check**:
   - Check for outdated dependencies
   - Check for known vulnerabilities
   - Generate SBOM

2. **Build**:
   - Build container images
   - Scan images for vulnerabilities
   - Sign images with Cosign

3. **Push**:
   - Push images to registry
   - Attach SBOM to images
   - Verify signatures

### Deployment Time

1. **Pre-Deployment**:
   - Verify image signatures
   - Check vulnerability scan results
   - Validate SBOM

2. **Deployment**:
   - Deploy only verified images
   - Log deployment details
   - Track image versions

### Runtime

1. **Monitoring**:
   - Monitor for new vulnerabilities
   - Alert on security issues
   - Track dependency updates

---

## Security Policies

### Image Security Policy

**Requirements**:
- ✅ Images must be signed
- ✅ Images must have SBOM
- ✅ Images must pass vulnerability scan
- ✅ Images must be from allowed registry

### Dependency Security Policy

**Requirements**:
- ✅ Dependencies must be up-to-date (< 30 days)
- ✅ Dependencies must not have critical vulnerabilities
- ✅ Dependencies must be from trusted sources

---

## Implementation Status

- [ ] SBOM generation automation
- [ ] Image signing setup (Cosign)
- [ ] Vulnerability scanning automation
- [ ] Dependency tracking setup
- [ ] Supply chain security CI/CD integration

---

## Referenties

- [SLSA Framework](https://slsa.dev/)
- [Sigstore](https://www.sigstore.dev/)
- [Cosign Documentation](https://docs.sigstore.dev/cosign/overview/)
- [Syft Documentation](https://github.com/anchore/syft)
- [Trivy Documentation](https://aquasecurity.github.io/trivy/)

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team, Security Team

