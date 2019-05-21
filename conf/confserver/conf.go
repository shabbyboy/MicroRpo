package confserver

import (
	"github.com/micro/go-config"
	"log"
)

type Conf struct{
	Path string
}
/*
这里用到了，就解释下把。如果是c 是指针代表可以在方法体内修改 Conf
而 c 是结构体表示不可修改Conf, 在方法里面做的修改，只在方法生命周期有效
 */
func (c Conf)LoacConf() error{
	err := config.LoadFile(c.Path)

	if err != nil{
		return err
	}
	return nil
}

func (c Conf)ConfExtract(cfg interface{}, keys ...string) error{
	err := config.Get(keys...).Scan(cfg)

	if err != nil{
		return err
	}
	return nil
}

func (c Conf)ConfChange(cfg interface{}, keys ...string){
	//开启一个监听，监听配置变化
	go func() {
		for {
			wcfg,err := config.Watch(keys...)
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


