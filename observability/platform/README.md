# Observability Platform

Dit directory bevat configuratie voor het observability platform.

## Overzicht

Het observability platform biedt unified observability voor het App Store platform met metrics, logs, traces en dashboards.

## Target Stack: LGTM

**Loki**: Log aggregation  
**Grafana**: Visualization  
**Tempo**: Distributed tracing  
**Mimir**: Metrics storage (Prometheus compatible)

**Alternative**: Grafana Cloud (managed LGTM stack)

## Components

### Metrics Collection

**Prometheus**:
- Metrics scraping
- Time-series storage
- Query language (PromQL)

**Mimir** (Future):
- Long-term metrics storage
- Scalable Prometheus backend
- Multi-tenancy support

**Location**: `observability/platform/metrics/`

---

### Log Collection

**Loki**:
- Log aggregation
- Label-based indexing
- Query language (LogQL)

**Agents**:
- Promtail (Loki agent)
- Kubernetes log collection
- Application log forwarding

**Location**: `observability/platform/logs/`

---

### Trace Collection

**Tempo**:
- Distributed tracing
- Trace storage
- Trace query language

**Instrumentation**:
- OpenTelemetry
- Application instrumentation
- Automatic trace collection

**Location**: `observability/platform/traces/`

---

### Visualization

**Grafana**:
- Unified dashboards
- Alerting
- Reporting

**Dashboards**:
- Single pane of glass
- Service-specific dashboards
- SLO dashboards

**Location**: `observability/dashboards/`

---

## Platform Architecture

### Collection Layer

**Metrics**:
- Prometheus exporters
- Kubernetes metrics API
- Application metrics endpoints

**Logs**:
- Loki agents (Promtail)
- Kubernetes log aggregation
- Application log forwarding

**Traces**:
- Tempo agents
- OpenTelemetry instrumentation
- Automatic trace collection

---

### Storage Layer

**Metrics**: Prometheus/Mimir  
**Logs**: Loki  
**Traces**: Tempo

---

### Query Layer

**Metrics**: PromQL  
**Logs**: LogQL  
**Traces**: TraceQL

---

### Visualization Layer

**Grafana**:
- Dashboards
- Alerts
- Reports

---

## Deployment

### Kubernetes Deployment

**Manifests**: `observability/platform/k8s/`

**Components**:
- Prometheus deployment
- Loki deployment
- Tempo deployment (future)
- Grafana deployment

### Configuration

**Prometheus**: `observability/platform/metrics/prometheus.yml`  
**Loki**: `observability/platform/logs/loki.yml`  
**Grafana**: `observability/platform/grafana/`

---

## Data Retention

**Metrics**: 30 days (Prometheus), 1 year (Mimir)  
**Logs**: 7 days (Loki)  
**Traces**: 7 days (Tempo)

---

## Referenties

- [Grafana LGTM Stack](https://grafana.com/docs/loki/latest/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [OpenTelemetry](https://opentelemetry.io/)

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team, SRE Team

