# 本地开发环境

docker-compose.yml 是本地开发环境，用于启动本地服务。

docker-devops-grafana-cloud.yml 是 grafana cloud 监控组件，用于监控本地服务。

docker-devops-local.yml 是本地监控组件，用于监控本地服务。(TODO)

## 启动本地服务

```shell
docker compose -f docker-compose.yml up -d
```

## 启动监控组件

```shell
docker compose -f docker-devops-grafana-cloud.yml up -d
```

```shell
docker compose -f docker-devops-local.yml up -d
```