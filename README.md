场景服务
====================

场景服务实现二维地图的九宫格拆分，管理玩家附近的广播、推送。


# 构建
在内网jenkins上选择项目`scene_go`，点击构建，构建成功后，会以docker镜像的方式上传到我们的docker私有仓库。

# 部署
1. 从私有仓库pull最新的镜像 `docker pull 139.198.2.55/soalib/scene:1.0_xxx`
2. 启动
```
# description
# MYSQL_SERVER mysql数据库地址
# MYSQL_AUTH mysql用户密码
# ZK_SERVERS zookeeper服务地址，可以是多个，用逗号分隔
# ZK_AUTH zookeeper用户密码
# ZK_REG_TCP 场景服务器tcp注册到zookeeper的地址，客户端用过该地址使用场景服务
# ZK_REG_HTTP 场景服务器http注册到zookeeper的地址

docker run -d \
        -e MYSQL_SERVER=127.0.0.1:3306 \
        -e MYSQL_AUTH=root:123456 \
        -e ZK_SERVERS=127.0.0.1:2181 \
        -e ZK_AUTH=platform:123456 \
        -e ZK_REG_TCP=127.0.0.1:9901 \
        -e ZK_REG_HTTP=127.0.0.1:9911 \
        --name scene \
        scene:build
```
3. 查看启动状态 `docker ps -a |grep scene`, 注意状态如果是`Exited`或不是`Up`证明启动失败
4. 查看日志 `docker logs scene`
