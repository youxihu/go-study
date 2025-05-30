package entity

type Config struct {
	DBFilePath struct {
		IP2RegionDBPath string `yaml:"ip2regionDBPath"`
		ASNDBPath       string `yaml:"asnDBPath"`
		CITYDBPath      string `yaml:"cityDBPath"`
	} `yaml:"dbfilepath"`

	LogFilesPath    []string `yaml:"logFilesPath"`
	DingTalkWebhook string   `yaml:"dingTalkWebhook"`
	WhiteList       []string `yaml:"whiteList"`
}
