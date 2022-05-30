package cmd

import (
	"bufio"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gproc"
	"github.com/linxlib/k/utils"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const platforms = `
    darwin    amd64
    darwin    arm64
    ios       amd64
    ios       arm64
    freebsd   386
    freebsd   amd64
    freebsd   arm
    linux     386
    linux     amd64
    linux     arm
    linux     arm64
    linux     ppc64
    linux     ppc64le
    linux     mips
    linux     mipsle
    linux     mips64
    linux     mips64le
    netbsd    386
    netbsd    amd64
    netbsd    arm
    openbsd   386
    openbsd   amd64
    openbsd   arm
    windows   386
    windows   amd64
	android   arm
	dragonfly amd64
	plan9     386
	plan9     amd64
	solaris   amd64
`

func GetMod(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(m, "module") {
			m = strings.TrimPrefix(m, "module")
			m = strings.TrimSpace(m)
			return m
		}
	}
	return ""
}

func Build() {
	g.Config().SetFileName("build.toml")
	parser, err := gcmd.Parse(g.MapStrBool{
		"n,name":    true,
		"v,version": true,
		"a,arch":    true,
		"s,system":  true,
		"p,path":    true,
		"t,tags":    true,
	})
	if err != nil {
		_log.Fatal(err)
	}
	file := parser.GetArg(2)
	if len(file) < 1 {
		// Check and use the main.go file.
		if utils.Exists("main.go") {
			file = "main.go"
		} else {
			_log.Fatal("编译文件不能为空")
		}
	}
	path := getOption(parser, "path", "./bin")
	name := getOption(parser, "name", gfile.Name(file))
	if len(name) < 1 || name == "*" {
		_log.Fatal("名称不能为空")
	}

	var (
		version       = getOption(parser, "version")
		archOption    = getOption(parser, "arch")
		systemOption  = getOption(parser, "system")
		tagsOption    = getOption(parser, "tags")
		customSystems = utils.SplitAndTrim(systemOption, ",")
		customArches  = utils.SplitAndTrim(archOption, ",")
	)

	if len(version) > 0 {
		path += "/" + version
	}
	// System and arch checks.
	var (
		spaceRegex  = regexp.MustCompile(`\s+`)
		platformMap = make(map[string]map[string]bool)
	)
	for _, line := range strings.Split(strings.TrimSpace(platforms), "\n") {
		line = utils.Trim(line)
		line = spaceRegex.ReplaceAllString(line, " ")
		var (
			array  = strings.Split(line, " ")
			system = strings.TrimSpace(array[0])
			arch   = strings.TrimSpace(array[1])
		)
		if platformMap[system] == nil {
			platformMap[system] = make(map[string]bool)
		}
		platformMap[system][arch] = true
	}
	modName := ""
	if utils.Exists("./go.mod") {
		modName = GetMod("go.mod")
	}

	ldFlags := ""

	// start building
	_log.Print("开始编译...")

	if utils.Exists("config.toml") {
		_log.Print("生成swagger...")
		if result, err := gproc.ShellExec("go run main.go -g"); err != nil {
			_log.Printf("生成失败, error:\n%s\n", utils.Trim(result))
		}
	} else {
		_log.Print("一般golang项目, 跳过swagger生成")
	}

	os.Setenv("CGO_ENABLED", "0")
	var (
		cmd  = ""
		ext  = ""
		tags = ""
	)
	for system, item := range platformMap {
		cmd = ""
		ext = ""
		if len(customSystems) > 0 && customSystems[0] != "all" && !utils.InArray(customSystems, system) {
			continue
		}
		for arch, _ := range item {
			if len(customArches) > 0 && customArches[0] != "all" && !utils.InArray(customArches, arch) {
				continue
			}
			if len(tagsOption) > 0 {
				tags = "-tags=" + tagsOption
			}
			if len(customSystems) == 0 && len(customArches) == 0 {
				if runtime.GOOS == "windows" {
					ext = ".exe"
				}
				ldFlags = fmt.Sprintf(`-X github.com/linxlib/kapi.VERSION=%s`, "NO_VERSION") +
					fmt.Sprintf(` -X github.com/linxlib/kapi.BUILDTIME=%s`, time.Now().Format("2006-01-02T15:04:01")) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.GOVERSION=%s`, runtime.Version()) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.OS=%s`, runtime.GOOS) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.ARCH=%s`, runtime.GOARCH) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.PACKAGENAME=%s`, modName)
				// Single binary building, output the binary to current working folder.
				output := "-o " + name + ext
				cmd = fmt.Sprintf(`go build %s %s -ldflags "%s"  %s`, tags, output, ldFlags, file)
			} else {
				// Cross-building, output the compiled binary to specified path.
				if system == "windows" {
					ext = ".exe"
				}
				os.Setenv("GOOS", system)
				os.Setenv("GOARCH", arch)
				ldFlags = fmt.Sprintf(`-X github.com/linxlib/kapi.VERSION=%s`, version) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.BUILDTIME=%s`, time.Now().Format("2006-01-02T15:04:01")) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.GOVERSION=%s`, runtime.Version()) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.OS=%s`, system) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.ARCH=%s`, arch) +
					fmt.Sprintf(` -X github.com/linxlib/kapi.PACKAGENAME=%s`, modName)
				cmd = fmt.Sprintf(
					`go build %s -o %s/%s/%s%s -ldflags "%s" %s`,
					tags, path, system+"_"+arch, name, ext, ldFlags, file,
				)
			}
			_log.Debug(cmd)
			// It's not necessary printing the complete command string.
			cmdShow, _ := utils.ReplaceString(`\s+(-ldflags ".+?")\s+`, " ", cmd)
			_log.Print(cmdShow)
			if result, err := gproc.ShellExec(cmd); err != nil {
				_log.Printf("编译失败, os:%s, arch:%s, error:\n%s\n", system, arch, utils.Trim(result))
			}
			if utils.Exists("gen.gob") {
				utils.CopyFile("gen.gob", fmt.Sprintf(
					`%s/%s/gen.gob`,
					path, system+"_"+arch))
				_log.Debug("拷贝gen.gob文件")
			}
			if utils.Exists("swagger.json") {
				utils.CopyFile("swagger.json", fmt.Sprintf(
					`%s/%s/swagger.json`,
					path, system+"_"+arch))
				_log.Debug("拷贝swagger.json文件")
			}
			if utils.Exists("config.toml") && !utils.Exists(fmt.Sprintf(
				`%s/%s/config.toml`,
				path, system+"_"+arch)) {
				utils.CopyFile("config.toml", fmt.Sprintf(
					`%s/%s/config.toml`,
					path, system+"_"+arch))
				_log.Debug("拷贝config.toml文件")
			}
			// single binary building.
			if len(customSystems) == 0 && len(customArches) == 0 {
				goto buildDone
			}
		}
	}
buildDone:
	_log.Print("完成!")

}

const nodeNameInConfigFile = "k"

// getOption retrieves option value from parser and configuration file.
// It returns the default value specified by parameter `value` is no value found.
func getOption(parser *gcmd.Parser, name string, value ...string) (result string) {
	result = parser.GetOpt(name)
	if result == "" && g.Config().Available() {
		result = g.Config().GetString(nodeNameInConfigFile + "." + name)
	}
	if result == "" && len(value) > 0 {
		result = value[0]
	}
	return
}
