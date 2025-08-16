load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_resource', 'helm_resource', 'helm_repo')

# Consul
helm_repo('hashicorp', 'https://helm.releases.hashicorp.com')
helm_resource(
    'consul',
    'hashicorp/consul',
    namespace='consul',
    flags=[
        '--namespace=consul',
        '--create-namespace',
        '--set=global.name=consul',
        '--values=./infra/helm/values/dev/_consul-values.yaml',
    ],
    pod_readiness='ignore',
    resource_deps=['hashicorp'],
    labels='tooling',
)
k8s_resource(
    'consul',
    port_forwards=['8501:8500'],
    labels='tooling',
    extra_pod_selectors=[{'component': 'server'}],
    discovery_strategy='selectors-only',
)

# API Gateway
gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway/cmd/main.go'
if os.name == 'nt':
    gateway_compile_cmd = './infra/docker/dev/api-gateway-build.bat'

local_resource(
    'api-gateway-compile',
    gateway_compile_cmd,
    deps=['./services/api-gateway', './shared'],
    labels='compiles',
)

docker_build_with_restart(
    'vasapolrittideah/moneylog-api-api-gateway',
    '.',
    entrypoint=['/app/build/api-gateway'],
    dockerfile='./infra/docker/dev/api-gateway.Dockerfile',
    only=['./build/api-gateway', './shared'],
    live_update=[
        sync('./build', '/app/build'),
        sync('./shared', '/app/shared')
    ],
)

k8s_yaml(helm(
    './infra/helm/charts/api-gateway',
    name='api-gateway',
    values=['./infra/helm/values/dev/api-gateway-values.yaml']
))

k8s_resource(
    'api-gateway',
    resource_deps=['api-gateway-compile', 'consul'],
    port_forwards='9000',
    labels='services',
)

# Auth Service
helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
helm_resource(
    'auth-mongodb',
    'bitnami/mongodb',
    flags=['--values=./infra/helm/values/dev/auth-mongodb-values.yaml'],
    resource_deps=['bitnami'],
)
k8s_resource('auth-mongodb', labels='databases')

auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/auth-service ./services/auth-service/cmd/main.go'
if os.name == 'nt':
    auth_compile_cmd = './infra/docker/dev/auth-service-build.bat'

local_resource(
    'auth-service-compile',
    auth_compile_cmd,
    deps=['./services/auth-service', './shared'],
    labels='compiles',
)

docker_build_with_restart(
    'vasapolrittideah/moneylog-api-auth-service',
    '.',
    entrypoint=['/app/build/auth-service'],
    dockerfile='./infra/docker/dev/auth-service.Dockerfile',
    only=['./build/auth-service', './shared'],
    live_update=[
        sync('./build', '/app/build'),
        sync('./shared', '/app/shared')
    ],
)

k8s_yaml(helm(
    './infra/helm/charts/auth-service',
    name='auth-service',
    values=[
        './infra/helm/values/dev/auth-service-values.yaml',
        './infra/helm/values/dev/auth-service-secrets.yaml',
    ]
))

k8s_resource(
    'auth-service',
    resource_deps=['auth-service-compile', 'auth-mongodb', 'consul'],
    labels='services',
)
