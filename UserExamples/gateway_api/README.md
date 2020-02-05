# API网关

## 使用

### 运行

```bash
> logs/access.log && go run main.go lib.go   --registry=etcd --registry_address=127.0.0.1:2379,127.0.0.1:2379   api  --handler=web --namespace=go.micro.tc   --address=0.0.0.0:8080
```

### 编译打包

打包

```
make build
```

运行二进制文件

```
./auth-srv
```

打包成docker镜像

```
make docker
```