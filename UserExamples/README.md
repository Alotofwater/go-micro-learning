# micro

## 实现功能
- 链路跟踪

```text
gateway-api -> api -> grpc
```
- 日志

```text
访问日志、业务日志: 增加 Trace-Id 字段
```
- etcd配置中心
```text
micro etcd配置库功能进行封装，动态更新配置增加HookFunc功能
```
- 其他主要功能
```text
熔断
prometheus监控
内部服务采用grpc通信
```

## 环境准备

- jaeger部署

```text
docker run -d --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.6
```

+ etcd部署

docker命令
```text
docker run -d \
  -p 2379:2379 \
  -p 2380:2380 \
  -v /Users/fred/application/data/etcd:/etcd-data \
  --name etcd_gcr_v3.3.13 \
  quay.io/coreos/etcd:v3.3.13 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new
```
客户端命令：
```text
export ETCDCTL_API=3
etcdctl --endpoints='127.0.0.1:2380' put  /gateway/config/log_config '{"log_file_name":"./logs/info.log","ac_file_name":"./logs/access.log","max_size":128,"max_backups":30,"max_a
ge":7,"log_level":"debug"}'

etcdctl --endpoints='127.0.0.1:2380' put  /gateway/config/redis_config '{"addr":"127.0.0.1:26379","password":"","timeout":2000,"dbNum":0,"enabled":true,"sentinel":{"enabled":fals
e,"master":"","nodes":""}}'

etcdctl --endpoints='127.0.0.1:2380' put  /gateway/config/db_config '{"driver":"mysql","url":"root:qwe123a@tcp(127.0.0.1:23306)/cmdb_dome?charset=utf8&parseTime=true&loc=Local","
enabled":true,"max_idle_connection":2,"max_open_connection":3}'
```  
- prometheus

配置文件
```text
global:
  scrape_interval: 15s
  scrape_timeout: 10s
  evaluation_interval: 15s
alerting:
  alertmanagers:
  - static_configs:
    - targets: []
    scheme: http
    timeout: 10s
scrape_configs:
  - job_name: APIGW
    honor_timestamps: true
    scrape_interval: 15s
    scrape_timeout: 10s
    metrics_path: /metrics
    scheme: http
    static_configs:
  - targets:
    - 10.104.34.106:8080
```
启动命令
```text
prometheus （启动时依赖本机配置文件 /tmp/conf.yml , 可更改命令自定义路径）
docker run --name prometheus -d -p 0.0.0.0:9090:9090 -v /Users/fred/application/data/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

```

