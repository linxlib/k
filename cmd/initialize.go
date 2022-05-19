package cmd

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"io/ioutil"
	"os"
	"os/exec"
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
	configFileContent = `# mysql 数据库配置
#[db]
#   enable = false 
#   mysql = "root:root@tcp(127.0.0.1:3306)/test?charset=utf8"

# redis 配置
#[redis]
#   enable = false
#   address = "127.0.0.1:6379"
#   password = ""
#   db = 1

[server]
    debug = true
    needDoc=true
    docName="KApi"
    docDesc="KApi Swagger Doc"
    port=2021
    openDocInBrowser=false
    docDomain=""
    docVer="v1"
    redirectToDocWhenAccessRoot=true
    apiBasePath="/"
    staticDirs=["static"]
    [server.cors]
        allowHeaders = ["Origin","Content-Length","Content-Type","Authorization","x-requested-with"]

`
	mainContent = `package main

import (
	"gitee.com/kirile/kapi"
)

func main() {
	k := kapi.New(func(option *kapi.Option) {
		//option.SetDocDomain("")
	})
	//k.RegisterRouter(new(api.CategoryController))

	k.Run()
}
`
)

func Initialize() {
	//TODO: 写出默认的配置文件
	if gfile.Exists("go.mod") {
		modName := GetMod("go.mod")
		if !gfile.Exists("build.toml") {
			r := fmt.Sprintf(buildFileContent, modName)
			ioutil.WriteFile("build.toml", []byte(r), os.ModePerm)
			_log.Println("写出build.toml")
		}
		if !gfile.Exists("config.toml") {
			//r := fmt.Sprintf(configFileContent,modName)
			ioutil.WriteFile("config.toml", []byte(configFileContent), os.ModePerm)
			_log.Println("写出config.toml")
		}
		if !gfile.Exists("api") {
			_log.Println("创建api目录")
			gfile.Mkdir("api")
		}
		if !gfile.Exists("main.go") {
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
	} else {
		_log.Println("go.mod不存在")
	}
}
