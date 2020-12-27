package conf

type model struct {
	Clickhouse struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Database string `yaml:"database"`
		Nums     int    `yaml:"nums"`
	}

	Debug    bool `yaml:"debug"`
	Interval int  `yaml:"interval"`
}

//Config 配置信息
var Config *model
