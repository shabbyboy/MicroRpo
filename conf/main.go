package main

/*
此main 函数 主要是开发时 测试用，并无其他用处，可以删除
 */

type Host struct {
	Address string `json:"address"`
	Port int `json:"port"`
	Host map[string]interface{} `json:"host"`
}

type Config struct {
	Hosts Host `json:"hosts"`
}

func main(){
	//conf := config.NewConfig()
	//这个路径和我常识不太一样

}
