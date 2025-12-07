# DORA Metrics Collection

Dit directory bevat de implementatie voor het verzamelen van DORA (DevOps Research and Assessment) metrics.

## DORA Four Key Metrics

### 1. Deployment Frequency
Hoe vaak wordt er succesvol gedeployed naar productie?

**Target**: > 10 deployments/dag (Elite)

**Collection**:
- Kubernetes deployment events
- Platform API deployment logs
- Git commit to deployment time

**Files**:
- ✅ `deployment-frequency/collector.go` - Metrics collector
- [ ] `deployment-frequency/dashboard.json` - Grafana dashboard config (in main dashboard)

### 2. Lead Time for Changes
Tijd van code commit tot deployment in productie.

**Target**: < 2 uur (Elite)

**Collection**:
- Git commit timestamps
- Deployment timestamps
- CI/CD pipeline duration

**Files**:
- ✅ `lead-time/collector.go` - Metrics collector
- [ ] `lead-time/dashboard.json` - Grafana dashboard config (in main dashboard)

### 3. Change Failure Rate
Percentage van deployments die falen in productie.

**Target**: < 5% (Elite)

**Collection**:
- Deployment status (success/failure)
- Rollback events
- Incident reports

**Files**:
- ✅ `change-failure-rate/collector.go` - Metrics collector
- [ ] `change-failure-rate/dashboard.json` - Grafana dashboard config (in main dashboard)

### 4. Mean Time to Recovery (MTTR)
Gemiddelde tijd om te herstellen van een failure.

**Target**: < 30 minuten (Elite)

**Collection**:
- Incident start/end times
- Deployment rollback times
- Recovery automation logs

**Files**:
- ✅ `mttr/collector.go` - Metrics collector
- [ ] `mttr/dashboard.json` - Grafana dashboard config (in main dashboard)

## Implementation Status

- ✅ Deployment frequency collector
- ✅ Lead time collector
- ✅ Change failure rate collector
- ✅ MTTR collector
- ✅ Grafana dashboard (`dashboard/dora-metrics.json`)
- ✅ Prometheus metrics exporter (`prometheus-exporter.go`)
- 🔄 Automated collection pipeline (basis geïmplementeerd)

## Data Sources

1. **Kubernetes API**: Deployment events, pod status
2. **Platform API**: Deployment logs, status endpoints
3. **Git**: Commit timestamps, branch information
4. **CI/CD**: Pipeline execution times
5. **Monitoring**: Incident tracking, alert resolution times

## Metrics Format

Alle metrics worden geëxporteerd in Prometheus format:

```prometheus
# Deployment frequency
platform_deployments_total{namespace="nl-appstore-registry",status="success"} 42

# Lead time
platform_lead_time_seconds{namespace="nl-appstore-registry",deployment="app-123"} 3600

# Change failure rate
platform_deployment_failures_total{namespace="nl-appstore-registry"} 2
platform_deployments_total{namespace="nl-appstore-registry"} 100

# MTTR
platform_recovery_time_seconds{namespace="nl-appstore-registry",incident="inc-456"} 1800
```

## Prometheus Exporter

De Prometheus exporter (`prometheus-exporter.go`) verzamelt metrics en exporteert ze via een `/metrics` endpoint.

**Gebruik**:
```bash
go run monitoring/dora-metrics/prometheus-exporter.go \
  -kubeconfig ~/.kube/config \
  -namespace nl-appstore-registry \
  -port 9090
```

**Metrics endpoint**: `http://localhost:9090/metrics`

## Dashboard

Grafana dashboard configuratie beschikbaar in:
- ✅ `dashboard/dora-metrics.json` - Complete DORA metrics dashboard

**Import in Grafana**:
1. Ga naar Grafana → Dashboards → Import
2. Upload `dashboard/dora-metrics.json`
3. Configureer Prometheus data source
4. Dashboard is klaar voor gebruik

## Next Steps

1. ✅ Collectors geïmplementeerd
2. ✅ Prometheus exporter gemaakt
3. ✅ Grafana dashboard gemaakt
4. 🔄 Integratie met CI/CD pipeline
5. 🔄 Automatische metrics collection in Kubernetes
6. 🔄 Alerting rules voor DORA metrics

---

**Status**: Basis implementatie voltooid  
**Eigenaar**: Platform Engineering Team
