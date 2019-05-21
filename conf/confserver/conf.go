package confserver

import (
	"github.com/micro/go-config"
	"log"
)

func ConfChange(cfg *interface{}, keys ...string){
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

func ConfExtract(cfg *interface{}, keys ...string){
	config.Get(keys...).Scan(cfg)
}
