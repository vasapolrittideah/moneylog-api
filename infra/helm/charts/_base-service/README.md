# base-service Helm Chart

This is a reusable base Helm chart (library chart) that provides common templates and configurations for microservices. It eliminates duplication while keeping individual service charts separate and maintainable.

## Features

- **Deployment**: Standard Kubernetes deployment with configurable replicas, resources, health checks, and environment variables
- **Service**: ClusterIP service with configurable ports and annotations
- **ConfigMap**: Optional configuration management with data injection
- **Secret**: Optional secret management with string data support
- **Security**: Pod and container security contexts with best practices
- **Image Management**: Flexible image configuration with global registry support
- **Environment Variables**: Support for both direct env vars and envFrom (ConfigMap/Secret references)

## Usage

### 1. Add as Dependency

In your service's `Chart.yaml`, add the base-service as a dependency:

```yaml
apiVersion: v2
name: your-service
type: application
version: 0.1.0
description: Your service chart
appVersion: "1.0.0"
dependencies:
  - name: base-service
    version: 0.1.0
    repository: file://../base-service
```

### 2. Create Templates

In your service's `templates/` directory, create simple template files that include the base templates:

**templates/deployment.yaml:**
```yaml
{{- include "base-service.deployment" . }}
```

**templates/service.yaml:**
```yaml
{{- include "base-service.service" . }}
```

**templates/configmap.yaml:**
```yaml
{{- include "base-service.configmap" . }}
```

**templates/secrets.yaml:**
```yaml
{{- include "base-service.secrets" . }}
```

### 3. Configure Values

In your service's `values.yaml` or `values-dev.yaml`, override the _base_service values:

```yaml
# Override base-service values for your service
nameOverride: "your-service"
fullnameOverride: "your-service"

# Global configuration (optional)
global:
  registry: "your-registry.com"
  imageTag: "v1.0.0"
  imagePullPolicy: IfNotPresent

# Deployment configuration
deployment:
  replicaCount: 2
  image:
    repository: your-service
    tag: "v1.0.0"  # Can be overridden by global.imageTag
    pullPolicy: IfNotPresent  # Can be overridden by global.imagePullPolicy
  
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 500m
      memory: 512Mi
  
  # Health checks
  livenessProbe:
    httpGet:
      path: /health
      port: http
    initialDelaySeconds: 30
    periodSeconds: 10
  
  readinessProbe:
    httpGet:
      path: /ready
      port: http
    initialDelaySeconds: 5
    periodSeconds: 5
  
  # Environment variables from ConfigMap/Secret
  envFrom: true

# Service configuration
service:
  type: ClusterIP
  port: 8080
  targetPort: 8080
  annotations: {}

# ConfigMap configuration
configMap:
  enabled: true
  data:
    ENVIRONMENT: "production"
    SERVICE_NAME: "your-service"
    SERVICE_PORT: "8080"

# Secret configuration
secrets:
  enabled: true
  stringData:
    DATABASE_PASSWORD: "your-secret-password"
    API_KEY: "your-api-key"

# Security contexts
podSecurityContext:
  fsGroup: 2000

securityContext:
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 1000

# Node selection
nodeSelector: {}
tolerations: []
affinity: {}
```

## Available Templates

- `base-service.deployment` - Kubernetes Deployment
- `base-service.service` - Kubernetes Service  
- `base-service.configmap` - Kubernetes ConfigMap
- `base-service.secrets` - Kubernetes Secret

## Helper Functions

- `base-service.name` - Chart name
- `base-service.fullname` - Full application name
- `base-service.chart` - Chart name and version
- `base-service.labels` - Common labels
- `base-service.selectorLabels` - Selector labels
- `base-service.image` - Full image name with registry and tag
- `base-service.envFrom` - Environment variables from ConfigMap/Secret helper

## Configuration Options

Key configuration sections available:

- `global.*` - Global settings (registry, imageTag, imagePullPolicy)
- `deployment.*` - Deployment configuration (replicas, image, resources, health checks, env)
- `service.*` - Service configuration (type, port, targetPort, annotations)
- `configMap.*` - ConfigMap configuration (enabled, data)
- `secrets.*` - Secret configuration (enabled, stringData)
- `podSecurityContext` - Pod-level security context
- `securityContext` - Container-level security context
- `nodeSelector` - Node selection constraints
- `tolerations` - Pod tolerations
- `affinity` - Pod affinity rules
- `imagePullSecrets` - Image pull secrets
- `podAnnotations` - Pod annotations

## Example Service Structure

```
your-service/
├── Chart.yaml          # Dependencies declaration
├── values-dev.yaml     # Development values
├── values-prod.yaml    # Production values
└── templates/
    ├── deployment.yaml # {{- include "base-service.deployment" . }}
    ├── service.yaml    # {{- include "base-service.service" . }}
    ├── configmap.yaml  # {{- include "base-service.configmap" . }}
    ├── secrets.yaml    # {{- include "base-service.secrets" . }}
    └── custom.yaml     # Any service-specific resources
```

This approach keeps your service charts DRY while maintaining flexibility for service-specific customizations.

## Image Configuration

The chart supports flexible image configuration:

1. **Global registry**: Set `global.registry` to prefix all images
2. **Global tag**: Set `global.imageTag` to override all image tags
3. **Per-service override**: Use `deployment.image.repository` and `deployment.image.tag`

Image resolution priority:
- `deployment.image.tag` > `global.imageTag` > `Chart.AppVersion`
- Final image: `{global.registry}/{deployment.image.repository}:{resolved-tag}`

## Environment Variables

Two approaches for environment variables:

1. **Direct variables** (not implemented in current templates):
```yaml
deployment:
  env:
    - name: SERVICE_NAME
      value: "my-service"
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-secret
          key: password
```

2. **From ConfigMap/Secret** (current implementation):
```yaml
deployment:
  envFrom: true  # Enables envFrom in deployment

configMap:
  enabled: true
  data:
    SERVICE_NAME: "my-service"

secrets:
  enabled: true
  stringData:
    DB_PASSWORD: "secret-password"
```

## Security Best Practices

The chart includes security best practices:
- Non-root user execution (runAsUser: 1000)
- Read-only root filesystem
- Dropped capabilities
- Security contexts for both pod and container levels

## Development Tips

1. **Testing locally**: Use `helm template` to validate your templates
2. **Debugging**: Use `--debug` flag with helm commands
3. **Values validation**: Always test with your actual values files
4. **Dependencies**: Run `helm dependency update` after adding the base-service dependency
