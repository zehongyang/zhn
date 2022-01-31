package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

var (
	gConfigData []byte
	once sync.Once
	err error
)

const (
	YmlSource = "ymlSource"
)

//加载yaml配置
func Load(out interface{}) error {
	once.Do(func() {
		err = loadData()
	})
	if err != nil {
		return err
	}
	if len(gConfigData) > 0 {
		return yaml.Unmarshal(gConfigData,out)
	}
	return nil
}

func loadData() error {
	ymls := os.Getenv(YmlSource)
	var loadFiles []string
	if len(ymls) < 1 {
		loadFiles = append(loadFiles,"./application.yml")
	}else{
		ss := strings.Split(ymls, ";")
		for _, s := range ss {
			loadFiles = append(loadFiles,s)
		}
	}
	for _, file := range loadFiles {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		gConfigData = append(gConfigData,b...)
	}
	return nil
}

type Flags struct {
	Debug bool
	Env string
}

func (s *Flags) Init ()  {
	flag.BoolVar(&s.Debug,"d",false,"DEBUG")
	flag.StringVar(&s.Env,"e","local.env","env file")
	testing.Init()
	flag.Parse()
}

func (s *Flags) IsDebug () bool {
	return s.Debug
}
