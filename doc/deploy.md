场景服务部署文档
=============================

# 0.依赖
应用程序 | 版本 | 说明
---|---|---
mysql|>=5.6|存放场景服务数据


------------------------------------------------------------

# 1.获取场景服务二进制包
前往jenkins( http://192.168.31.202:8000/job/scene )构建scene，构建成功后下载最新的`scene_xxx.tar.gz`包


------------------------------------------------------------

# 2.更改配置

#### 1. 解压`scene_xxx.tar.gz`，得到`scene`目录

#### 2. 更改启动节点名

```
//打开配置文件
vi scene/release/{VERSION}/vm.args

//更改节点名(如下第二行)，搭建集群时才有必要改，默认scene_master1@127.0.0.1
## Name of the node
-name scene_master1@127.0.0.1
```

#### 3. 更改启动参数

```
//打开配置文件
vi scene/release/{VERSION}/sys.config


//1) 更改master节点名(如下第三行)，搭建集群时才有必要改，默认['scene_master1@127.0.0.1']
    {cluster, [
        {role, master},
        {masters, ['scene_master1@127.0.0.1']},
        {tables, [session, app_proc, space_proc, grid_proc, super]}
    ]},


//2) 更改场景服务注册地址(如下第三行)，改为场景服务对外提供服务的ip
    {scene, [
        {mode, release},
        {zk_register_ip, "temp_zk_register_ip"},
        {tcp, [{port,9901}, {max_ac,10}, {max_conn,10240}, {type,1}]},
        {http, [{port,9911}, {max_ac,100}]}
    ]},


//3) 更改zookeeper地址(如下第二行)，和账号密码(如下第三行)
    {framework, [
        {zk, {"172.16.1.23", 2181}},
        {zk_account, {"platform", "jhqc_zk@#$%2016"}}
    ]},


//4) 更改mysql地址(如下第五行)和端口(如下第六行)
        {mysql_poolboy, [
            {pool1, {
                [{size, 10}, {max_overflow, 20}],
                [
                    {host, "172.16.1.27"},
                    {port, 3307},
                    {user, "root"},
                    {password, "123456"},
                    {database, "scene_db"}
                ]}}
        ]}
```


------------------------------------------------------------

# 3.准备数据库

执行场景服务的sql脚本：`scene/db/db.sql`


------------------------------------------------------------

# 4.启动/停止场景服务

```
//启动
chmod +x scene/bin/scene
scene/bin/scene start
//停止
scene/bin/scene stop
```


------------------------------------------------------------

# 5.检查是否已启动

#### 1. 查看场景服务进场是否还在

```
ps -ef |grep scene
```

#### 2.查看启动日志

```
cat scene/log/erlang.log.XXX

//如果启动正常，输出日志的最后内容大致如下：
01:57:25.188 [info] Application super started on node 'scene_master1@192.168.31.216'
01:57:25.188 [info] Application poolboy started on node 'scene_master1@192.168.31.216'
01:57:25.188 [info] Application mysql started on node 'scene_master1@192.168.31.216'
01:57:25.193 [info] Application mysql_poolboy started on node 'scene_master1@192.168.31.216'
01:57:25.193 [info] Application reloader started on node 'scene_master1@192.168.31.216'
01:57:25.193 [info] Application cowlib started on node 'scene_master1@192.168.31.216'
01:57:25.193 [info] Application ranch started on node 'scene_master1@192.168.31.216'

------------------"src/scene_app.erl":13-----------------
------------------------scene start--------------------------

------------------"src/scene.erl":34-----------------
--------------------master-------------------
ranch opts [{port,9901},{max_ac,10},{max_conn,10240},{type,1}]
cowboy opts [{port,9911},{max_ac,100}]
01:57:25.198 [info] register service path <<"/platform/scene/tcp/scene_master1@192.168.31.216">> data <<"{\"uri\":\"tcp://192.168.31.216:9901\"}">>
01:57:25.198 [info] register zk tcp-uri = <<"tcp://192.168.31.216:9901">>, ret = ok
01:57:25.201 [info] register service path <<"/platform/scene/http/scene_master1@192.168.31.216">> data <<"{\"uri\":\"http://192.168.31.216:9911\"}">>
01:57:25.202 [info] register zk http-uri = <<"http://192.168.31.216:9911">>, ret = ok

------------------"src/scene.erl":26-----------------
--------------------start-------------------
cookie = jhqc_platform_scene
Eshell V7.2  (abort with ^G)
```

#### 3.查看zookeeper是否已注册成功

**注意**：ip换成 `2-3-3)章节` 中的zookeeper地址

http://192.168.31.216:8011/get_service?path=/platform/scene/http

http://192.168.31.216:8011/get_service?path=/platform/scene/tcp

