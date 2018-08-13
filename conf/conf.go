/*
Package conf 用于项目基本配置。
*/
package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/simplejia/utils"
)

// Conf 定义配置参数
type Conf struct {
	// 基本配置
	App *struct {
		Name string
		Port int
	}
	// clog日志输出配置
	Clog *struct {
		Name  string
		Mode  int
		Level int
	}
	// 各种名字或addr配置
	Addrs *struct {
		Clog string
	}
}

var (
	// Env 代表当前运行环境
	Env string
	// C 代表当前运行配置对象
	C *Conf
	// TestCase 运行单元测试时，会设为true
	TestCase bool
)

func init() {
	var env string
	flag.StringVar(&env, "env", "prod", "set env")
	var test bool
	flag.BoolVar(&test, "test", false, "set test case flag")
	flag.Parse()

	Env = env
	TestCase = test

	dir := "conf"
	for i := 0; i < 3; i++ {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			break
		}
		dir = filepath.Join("..", dir)
	}
	fcontent, err := ioutil.ReadFile(filepath.Join(dir, "conf.json"))
	if err != nil {
		log.Printf("get conf file contents error: %v\n", err)
		os.Exit(-1)
	}

	fcontent = utils.RemoveAnnotation(fcontent)
	var envs map[string]*Conf
	if err := json.Unmarshal(fcontent, &envs); err != nil {
		log.Printf("conf.json wrong format: %v\n", err)
		os.Exit(-1)
	}

	C = envs[env]
	if C == nil {
		log.Printf("env not right: %s\n", env)
		os.Exit(-1)
	}

	log.Printf("env: %s\nconf: %s\n", env, utils.Iprint(C))
}
