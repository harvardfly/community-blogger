# community-blogger
gin+grpc+wire的脚手架，集成MySQL、Redis、ES、Kafka、分布式限流、服务注册发现、jaeger链路追踪、prometheus监控、cronlib定时调度任务及支付宝支付API，可作为中型互联网公司golang架构实践

## 技术选型
```$xslt
1. web架构：gin+gorm+redis
2. 数据库：mysql
3. logger日志zap  日志切割lumberjack
4. 依赖注入wire --> 解决大量配置文件初始化问题
5. grpc(用户模块 调用rpc服务根据用户ID获取用户信息)
6. JWT认证
7. 基于redis实现分布式限流中间件（分别实现漏桶和令牌桶算法） 可用于防止用户恶意发布文章  --> 分布式限流 服务保护
8. etcd作为服务注册发现  --> 服务注册发现  负载均衡
9. opentracing+jaeger+ElasticSearch分布式链路追踪日志存储  --> 复杂调用链路问题排查
10. prometheus 监控服务指标  分析gc等问题
11. kafka作为消息中间件 解耦应用 --> kafka高吞吐率将极大改善大流量下服务器压力
12. ES作为日志存储和搜索服务 --> 发挥ES倒排索引优势，可用于文章的搜索
13. cron调度任务模块，支持任务注册、删除等  同步redis中文章浏览次数到ES中
14. redis实现并发锁，lua脚本原子性释放锁
```
## 功能模块 以博客系统作为脚手架示例
```$xslt
各模块单独抽出一个service 采用微服务架构分布式部署
考虑到首页访问流量过大的问题，将home单独抽成一个服务，方便后续负载均衡扩展
1. home(首页模块)
2. article(文章模块)
   文章模块有article和category
3. user模块 
    userRPC(用户RPC服务模块 grpc实现)
    userAPI模块 调用userRPC获取用户信息
4. redis实现文章浏览次数和热点文章排行榜(zset+goroutine实现 提高并发效率)
5. 认证中间件采用JWT
6. 单元测试(主要针对repositories数据库操作的测试)
```

## 项目部署方式
```$xslt
1. 采用docker swarm + hub 的方式部署(目前实现的)
2. 采用k8s + hub 的方式部署(预留)
```

## 项目目录结构
```$xslt
api         --      存放protobuf相关文件 供client和server通信调用
cmd         --      可执行程序的入口(可以有多个可执行程序，每一个的main函数都在子文件夹)
configs     --      项目配置文件(可以对应多个可执行程序有多个配置文件)
internal    --      包含app和pkg
            app         --      项目的逻辑代码，包含controllers repositories services
            pkg         --      通用的代码，项目的公共代码
vendor      --      项目依赖包/库
.gitignore  --      git忽略文件列表
.dockerignore  --   docker忽略文件列表
Dockerfile    --    docker镜像配置文件
docker-compose.yml    --    docker-compose配置文件（已gitignore 直接在protainer中配置）
go.mod      --      项目依赖的第三方包
go.sum      --      所有依赖的包
Makefile    --      执行脚本的Makefile文件
README.md   --      说明文档
```

## 系统环境要求
```$xslt
golang >= 1.13
```

## 项目本地启动执行
```$xslt
1. cd community-blogger
2. cp /configs/home.yaml.example home.yaml
3. cp /configs/article.yaml.example article.yaml
4. ./dist/manage -f configs/home.yaml
5. ./dist/manage -f configs/article.yaml
6.  与上述步骤类似,此处省略 ...
```

## golint 代码规范检查
```$xslt
1. cd community-blogger
2. go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status
```

## 单元测试示例
```$xslt
cd community-blogger/internal/app/article/repositories
go test -v
go test -cover
```

## 接口postman测试示例
```$xslt
调用rpc服务获取用户信息 JWT认证
GET http://127.0.0.1:8004/api/v1/user?id=1
herders: Authorization eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTc4MjYwMzR9.BrcFgcv8GMYxDlR7QU3U0EDd9UbmaSVwaKKOaydHgus
{
    "data": {
        "id": 1,
        "username": "aaa3",
        "token": "553f5acf-d370-4bda-9ccd-ac8eb64e9665"
    }
}
```

## CSRF TOKEN
```$xslt
CSRF TOKEN 认证
GET http://127.0.0.1:8001/api/v1/protected
{
    "csrf_token": "pZsqS3IJf5J-fejsO6ikdR6aXyw="
}

带csrf_token请求：
headers  key：X-CSRF-TOKEN value：pZsqS3IJf5J-fejsO6ikdR6aXyw=
```
