package cfg

import "github.com/spf13/viper"

type Cfg struct {
	Mysql   MsCfg
	Logfile string `json:"logfile"`
}

type MsCfg struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Driver   string `json:"driver"`
	Dbname   string `json:"dbname"`
}

func LoadConfig(cfg *Cfg) error {
	viper.AddConfigPath("cfg")
	viper.SetConfigName("cfg")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return err
	}

	return nil
}
