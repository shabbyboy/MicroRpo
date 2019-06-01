package confserver

import (
	"errors"
	"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-config/encoder"
	"github.com/micro/go-config/encoder/json"
	"github.com/micro/go-config/encoder/toml"
	"github.com/micro/go-config/encoder/yaml"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-config/source/file"
	"log"
)

type SourceType string
type EncoderType string

const (
	CONFIG_CONSUL SourceType = "consul"
	CONFIG_ETCD SourceType = "etcd"
	CONFIG_FILE SourceType = "file"
)

const(
	CODER_JSON EncoderType = "json"
	CODER_YAML EncoderType = "yaml"
	CODER_TOML EncoderType = "toml"
)


type Conf struct{
	SrcType SourceType
	Coder EncoderType
	Host string
	Path string
	newConfig config.Config
	//推出监视
	Exit chan error
}


var CoderMap = map[string]encoder.Encoder{
"toml":toml.NewEncoder(),
"json":json.NewEncoder(),
"yaml":yaml.NewEncoder(),
}

func DefaultConf(path string) Conf{
	conf := Conf{
		SrcType:CONFIG_FILE,
		Coder:CODER_JSON,
		Path:path,
	}
	conf.LoadConf()
	return conf
}

/*
如果是c 是指针代表可以在方法体内修改 Conf
而 c 是结构体表示不可修改Conf, 在方法里面做的修改，只在方法生命周期有效
 */
func (c *Conf)LoadConf() error{
	c.newConfig = config.NewConfig()

	var coder encoder.Encoder
	switch c.Coder {
	case CODER_JSON:
		coder = json.NewEncoder()
	case CODER_TOML:
		coder = toml.NewEncoder()
	case CODER_YAML:
		coder = yaml.NewEncoder()
	default:
		coder = json.NewEncoder()
	}

	switch c.SrcType {
	case CONFIG_FILE:
		src := file.NewSource(
			file.WithPath(c.Path),
			source.WithEncoder(coder),
			)
		c.newConfig.Load(src)
	case CONFIG_CONSUL:
		src := consul.NewSource(
			consul.WithAddress(c.Host),
			consul.WithPrefix(c.Path),
			consul.StripPrefix(true),
			)
		c.newConfig.Load(src)
	case CONFIG_ETCD:
		src := consul.NewSource(
			consul.WithAddress(c.Host),
			consul.WithPrefix(c.Path),
			consul.StripPrefix(true),
		)
		c.newConfig.Load(src)
	default:
		return errors.New("srcType is must or srctype is not optionally")
	}

	return nil
}

func (c *Conf)ConfExtract(cfg interface{}, keys ...string) error{
	err := c.newConfig.Get(keys...).Scan(cfg)

	if err != nil{
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *Conf)ConfChange(cfg interface{}, keys ...string){
	//开启一个监听，监听配置变化
	go func() {
		for {
			wcfg,err := c.newConfig.Watch(keys...)
			if err != nil{
				log.Println(err)
			}

			v,errn := wcfg.Next()

			if errn != nil{
				log.Println(errn)
			}

			errs := v.Scan(cfg)

			if errs != nil{
				log.Println(errs)
			}
		}
	}()
}


