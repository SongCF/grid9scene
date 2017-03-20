#!/bin/bash

rm -f docker/sys.config
cp docker/temp_sys.config docker/sys.config


# cluster
sed -i "s/temp_cluster_app_id/$CLUSTER_APP_ID/g" docker/sys.config
sed -i "s/temp_cluster_selection_host/$CLUSTER_SELECTION_HOST/g" docker/sys.config

# zk 获取该节点的公网ip
SCENE_ZK_REGISTER_IP=$(curl http://icanhazip.com)
sed -i "s/temp_scene_zk_register_ip/$SCENE_ZK_REGISTER_IP/g" docker/sys.config

# framework
sed -i "s/temp_framework_zk_ip/$FRAMEWORK_ZK_IP/g" docker/sys.config
sed -i "s/temp_framework_zk_port/$FRAMEWORK_ZK_PORT/g" docker/sys.config
sed -i "s/temp_framework_zk_account/$FRAMEWORK_ZK_ACCOUNT/g" docker/sys.config
sed -i "s/temp_framework_zk_password/$FRAMEWORK_ZK_PASSWORD/g" docker/sys.config

# mysql
sed -i "s/temp_mysql_host/$MYSQL_HOST/g" docker/sys.config
sed -i "s/temp_mysql_port/$MYSQL_PORT/g" docker/sys.config
sed -i "s/temp_mysql_account/$MYSQL_ACCOUNT/g" docker/sys.config
sed -i "s/temp_mysql_password/$MYSQL_PASSWORD/g" docker/sys.config

# lager
sed -i "s/temp_console_log_level/$LOG_LEVEL/g" docker/sys.config

# 更改erlang node ip
# 按照自己的网络更改 |grep 10.0|grep /24
# ip a|grep inet|grep 10.0|grep /24|awk '{print $2}'|awk -F"/" '{print $1}'
NODE_IP=$(ip a|grep inet|grep ${SUBNET}|grep ${SUBNET_MASK}|awk '{print $2}'|awk -F"/" '{print $1}')


erl -noshell \
    -name scene@${NODE_IP} \
    -setcookie jhqc_platform_scene \
    -pa ebin \
    -pa deps/*/ebin \
    -config docker/sys.config \
    -s scene