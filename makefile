
GOPACK = "./gp"

# 所有
all: deps
	go build -o ./scene

# 获取依赖
deps:
	$(GOPACK) get-deps


# 清除
clean:
	rm -rf vendor;rm -rf .gopack


# 生成协议.
pb:
	proto4go.exe -i ./pb -o ./pb


# 启动
s:
	./scene


# 单元测试


# 构建数据库
db:
	/bin/bash db/db.sh



# 构建docker镜像
docker:
	/bin/bash ./scripts/docker_build.sh




.PHONY:deps clean pb s db docker


