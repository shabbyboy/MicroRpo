package main

/*
go-config demo
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
