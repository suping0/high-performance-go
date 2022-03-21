package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	Name    string `json:"server-name"`
	IP      string `json:"server-ip"`
	URL     string `json:"server-url"`
	Timeout string `json:"timeout"`
}

func readConfig() *Config {
	// read from xxx.json，省略
	config := Config{}
	// 结构体的元数据信息(字段等信息)
	typ := reflect.TypeOf(config)
	// 返回指针所指向结构体，实际所存储的值
	value := reflect.Indirect(reflect.ValueOf(&config))
	// NumField() 返回struct的字段数量
	for i := 0; i < typ.NumField(); i++ {
		// 结构体字段
		f := typ.Field(i)
		// 获取字段的tag
		if v, ok := f.Tag.Lookup("json"); ok {
			key := fmt.Sprintf("CONFIG_%s", strings.ReplaceAll(strings.ToUpper(v), "-", "_"))
			// 从环境变量中读取值
			if env, exist := os.LookupEnv(key); exist {
				value.FieldByName(f.Name).Set(reflect.ValueOf(env))
			}
		}
	}
	return &config
}

func main() {
	os.Setenv("CONFIG_SERVER_NAME", "global_server")
	os.Setenv("CONFIG_SERVER_IP", "10.0.0.1")
	os.Setenv("CONFIG_SERVER_URL", "geektutu.com")
	c := readConfig()
	fmt.Printf("%+v", c)
}
