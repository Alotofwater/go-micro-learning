# 用户web服务

## 使用

### 运行

```bash
> logs/info.log && >logs/access.log &&  go run main.go lib.go 
```

### 编译打包

打包

```
make build
```

运行二进制文件

```
./user-web
```

打包成docker镜像

```
make docker
```