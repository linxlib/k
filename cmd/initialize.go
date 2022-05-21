package cmd

import (
	"fmt"
	"github.com/linxlib/k/utils"
	"github.com/linxlib/kapi"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	buildFileContent = `# 编译配置
[k]
   name = "%s" # 编译的可执行文件名
   version = "1.0.0" # 版本
   arch = "amd64"  # 平台
   system = "windows,linux" # 系统
   path = "./bin" # 输出目录
`

	mainContent = `package main

import (
	"github.com/linxlib/kapi"
)

func main() {
	k := kapi.New(func(option *kapi.Option) {

	})
	//k.RegisterRouter(new(api.CategoryController))

	k.Run()
}
`
	dockerFileContent = `
FROM golang:1.18 as build
MAINTAINER "yourname <youremail>"

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn" \
	GOPRIVATE="gitee.com"

RUN mkdir /src
RUN go install github.com/linxlib/k@latest
WORKDIR /src

COPY . .
RUN go mod tidy
RUN k build

FROM ubuntu as prod
RUN mkdir /app
WORKDIR /app
RUN export DEBIAN_FRONTEND=noninteractive  \
    && apt-get update \
    && apt-get install -y tzdata \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl \
    && dpkg-reconfigure --frontend noninteractive tzdata
COPY --from=build /src/bin/1.0.0/linux_amd64/<appname> .
COPY --from=build /src/bin/1.0.0/linux_amd64/gen.gob .
# COPY --from=build /src/bin/1.0.0/linux_amd64/config.toml  .
COPY --from=build /src/bin/1.0.0/linux_amd64/swagger.json .
RUN ln -fs /app/<appname> /usr/bin/<appname>
RUN apt-get clean
EXPOSE 1509

CMD ["<appname>"]



`
)

func Initialize() {
	//TODO: 写出默认的配置文件
	if utils.Exists("go.mod") {
		modName := GetMod("go.mod")
		if !utils.Exists("build.toml") {
			r := fmt.Sprintf(buildFileContent, modName)
			ioutil.WriteFile("build.toml", []byte(r), os.ModePerm)
			_log.Println("写出build.toml")
		}
		if !utils.Exists("config.toml") {
			//r := fmt.Sprintf(configFileContent,modName)
			ioutil.WriteFile("config.toml", kapi.GetEmptyConfig(), os.ModePerm)

			_log.Println("写出config.toml")
		}
		if !utils.Exists("api") {
			_log.Println("创建api目录")
			utils.Mkdir("api")
		}
		if !utils.Exists("main.go") {
			//r := fmt.Sprintf(mainContent, modName, modName)
			ioutil.WriteFile("main.go", []byte(mainContent), os.ModePerm)
			_log.Println("写出main.go")
			output, err := exec.Command("gofmt", "-l", "-w", "./").Output()
			if err != nil {
				_log.Error(err)
				return
			}
			_log.Println(string(output))
			output, err = exec.Command("go", "mod", "tidy").Output()
			if err != nil {
				_log.Error(err)
				return
			}

			_log.Println(string(output))
		}
		if !utils.Exists("Dockerfile") {
			a := strings.ReplaceAll(dockerFileContent, "<appname>", modName)
			ioutil.WriteFile("Dockerfile", []byte(a), os.ModePerm)
			_log.Println("写出Dockerfile")
		}
	} else {
		_log.Println("go.mod不存在")
	}
}
