
# 不要在jenkins配置中直接调用该脚本，因为：脚本内的错误不会导致构建失败

git checkout .
chmod +x rebar
chmod +x jenkins/test/test.sh
chmod +x jenkins/build.sh

# mysql
LOCAL_IP=$(ip a|grep inet|grep 192.168.|grep /24|awk '{print $2}'|awk -F"/" '{print $1}')
MYSQL_IP=$LOCAL_IP
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PS=root123
# rebuild db
mysql -u$MYSQL_USER -p$MYSQL_PS --default-character-set=utf8 -e "source db/db.sql"


docker run --rm \
        -v $(pwd):/jh/source/scene/ \
        -v  ~/.ssh:/root/.ssh \
        -w /jh/source/scene/ \
        -e MYSQL_IP=$MYSQL_IP \
        -e MYSQL_PORT=$MYSQL_PORT \
        -e MYSQL_USER=$MYSQL_USER \
        -e MYSQL_PS=$MYSQL_PS \
        -e ZOOKEEPER_IP=$LOCAL_IP \
        -e ZOOKEEPER_PORT=2181 \
        139.198.2.55:5000/elixir_mix \
        /bin/bash -c '/jh/source/scene/jenkins/test/test.sh'



# 打包使用docker镜像
# docker run --rm \
#         -v $(pwd):/jh/source/scene/ \
#         -w /jh/source/scene/ \
#         erlang:18 \
#         /bin/bash jenkins/build.sh
# 打包不使用docker镜像
make clean
make
make g
make
make rel



gitcommit=`git log | head -1 | awk '{print substr($2,0,6)}'`
datetime=`date "+%Y%m%d%H%M%S"`
dockertag="1.0.0_"$datetime"_"$gitcommit
cp _rel/scene_rel.tar.gz  scene_$dockertag.tar.gz