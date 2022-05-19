# KAPI Cli

## 安装

```shell
go install gitee.com/kirile/k-cli/k@latest
```
或者clone代码进行编译
```shell
git clone https://gitee.com/kirile/k-cli
cd k-cli/k
go build -o k.exe main.go
start k.exe
```
运行编译的可执行程序, 可以选择安装到某个PATH目录下

## 使用

```shell
k init
k build
k build -t=nomsgpack
k run
etc....
```

#### 命令

- `k init` 用于在`go.mod`存在的情况下为你生成`build.toml` `config.toml`  `main.go` 并执行`go mod tidy`

- `k run` 运行, 监听源码变动, 重新运行
- `k` 不带参数运行则默认为安装可执行程序到系统 PATH目录中
- `k build` 依赖项目根目录下的 build.toml 文件, **目前版本会在构建后一起拷贝 gen.gob  swagger.json 和配置文件到构建目录(配置文件已存在则不覆盖)**
```toml
# 编译配置
[k]
  name = "api_base" # 编译的可执行文件名
  version = "1.0.0" # 版本
  arch = "amd64"  # 平台
  system = "darwin" # 系统
  path = "./bin" # 输出目录
```
使用k编译的kapi程序, 运行时将输出相关日志信息, 例如
```
--------------------------------------------
    _/    _/    _/_/    _/_/_/    _/_/_/
   _/  _/    _/    _/  _/    _/    _/
  _/_/      _/_/_/_/  _/_/_/      _/
 _/  _/    _/    _/  _/          _/
_/    _/  _/    _/  _/        _/_/_/

 Version:   1.0.0/go1.16.6
 OS/Arch:   windows/amd64
 BuiltTime: 2021-08-20T11:31:08
-------------------------------------------- 
```
如果build.toml不存在, 则只会在./bin目录下输出当前系统的可执行程序(相当于 go build main.go), 可以输入`k init` 生成这个文件

- 新增 `k build -t=nomsgpack` 可以带tags进行编译
- 新增生成模式, 在编译前先生成gen.gob和swagger.json