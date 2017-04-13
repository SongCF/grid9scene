#!/bin/bash

# set scene-node host
# set zookeeper host
# set mysql host
SYS_CONFIG=./conf.ini
sed -i "s/\(db_server = \).*/\1$MYSQL_SERVER/g" $SYS_CONFIG
sed -i "s/\(db_auth = \).*/\1$MYSQL_AUTH/g" $SYS_CONFIG
sed -i "s/\(zk_servers = \).*/\1$ZK_SERVERS/g" $SYS_CONFIG
sed -i "s/\(zk_auth = \).*/\1$ZK_AUTH/g" $SYS_CONFIG
sed -i "s/\(zk_reg_tcp_addr = \).*/\1$ZK_REG_HTTP/g" $SYS_CONFIG
sed -i "s/\(zk_reg_http_addr = \).*/\1$ZK_REG_HTTP/g" $SYS_CONFIG

./scene
