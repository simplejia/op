/*
Package mongo 用于mongo连接定义，所有mongo db配置需在此目录定义，一个db对应一个文件。
*/
package mongo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/mgo.v2"

	"github.com/simplejia/op/conf"
	"github.com/simplejia/utils"
)

// Conf 用于mongo连接配置
type Conf struct {
	Dsn string // 连接串: mongo://127.0.0.1:27017...
}

var (
	// DBS 表示mongo连接，key是db名，value是db连接
	DBS map[string]*mgo.Session = map[string]*mgo.Session{}
)

func init() {
	dir := "mongo"
	for i := 0; i < 3; i++ {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			break
		}
		dir = filepath.Join("..", dir)
	}
	err := filepath.Walk(
		dir,
		func(path string, info os.FileInfo, err error) (reterr error) {
			if err != nil {
				reterr = err
				return
			}
			if info.IsDir() {
				return
			}
			if strings.HasPrefix(path, ".") {
				return
			}
			if filepath.Ext(path) != ".json" {
				return
			}

			fcontent, err := ioutil.ReadFile(path)
			if err != nil {
				reterr = err
				return
			}
			fcontent = utils.RemoveAnnotation(fcontent)
			var envs map[string]*Conf
			if err := json.Unmarshal(fcontent, &envs); err != nil {
				reterr = err
				return
			}

			c := envs[conf.Env]
			if c == nil {
				reterr = fmt.Errorf("env not right: %s", conf.Env)
				return
			}

			session, err := mgo.Dial(c.Dsn)
			if err != nil {
				reterr = err
				return
			}
			// 在跑单元测试时，避免db主从同步延迟导致读不到最新数据的问题
			if conf.TestCase {
				session.SetMode(mgo.PrimaryPreferred, true)
			} else {
				session.SetMode(mgo.SecondaryPreferred, true)
			}

			key := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			DBS[key] = session
			return
		},
	)
	if err != nil {
		log.Printf("conf(mongo) not right: %v\n", err)
		os.Exit(-1)
	}

	log.Printf("conf(mongo): %v\n", DBS)
}
