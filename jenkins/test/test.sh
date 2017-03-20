echo "scene: deps.get ..."
mix deps.get
mix deps.update --all
mix deps.compile

echo "scene: replace mysql ip ..."
sed -i "s/tmp_mysql_ip/$MYSQL_IP/g" config/config_jenkins.exs
sed -i "s/tmp_mysql_port/$MYSQL_PORT/g" config/config_jenkins.exs
sed -i "s/tmp_mysql_user/$MYSQL_USER/g" config/config_jenkins.exs
sed -i "s/tmp_mysql_ps/$MYSQL_PS/g" config/config_jenkins.exs


sed -i "s/tmp_zookeeper_ip/$ZOOKEEPER_IP/g" config/config_jenkins.exs
sed -i "s/'tmp_zookeeper_port'/$ZOOKEEPER_PORT/g" config/config_jenkins.exs


sed -i "s/tmp_leaderselection_ip/$LEADER_IP/g" config/config_jenkins.exs
sed -i "s/tmp_leaderselection_port/$LEADER_PORT/g" config/config_jenkins.exs


rm -f config/config.exs
cp config/config_jenkins.exs config/config.exs

echo "scene: begin test ..."
mix test
