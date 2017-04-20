GOPACK = "./gopack"

ifeq ("$(shell uname -o)", "Cygwin")
  GOPACK = "./gopack-win"
endif


# 获取依赖，然后编译
all: deps
	go build -o ./scene

# 获取依赖
deps:
	$(GOPACK) get-deps

fmt:
	gofmt -w .


# 清除
clean-deps:
	rm -rf vendor;rm -rf .gopack

clean:
	rm -rf ./scene


# 生成协议.
pb:
	proto4go.exe -i ./pb -o ./pb


# 启动
s:
	./scene


# 单元测试


# 构建数据库
db:
	/bin/bash scripts/db.sh



# 构建docker镜像
docker:
	/bin/bash ./scripts/docker_build.sh




.PHONY:deps clean clean-deps pb s db docker


