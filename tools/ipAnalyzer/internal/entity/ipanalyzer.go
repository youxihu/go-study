package entity

type Config struct {
	DBFilePath      DBFilePath      `yaml:"dbFilepath"`
	LogFilesPath    []string        `yaml:"logFilesPath"`
	DingTalkWebhook string          `yaml:"dingTalkWebhook"`
	WhiteList       []string        `yaml:"whiteList"`
	Thresholds      ThresholdConfig `yaml:"thresholds"`
}

type DBFilePath struct {
	IP2RegionDBPath string `yaml:"ip2regionDBPath"`
	ASNDBPath       string `yaml:"asnDBPath"`
	CITYDBPath      string `yaml:"cityDBPath"`
}

type ThresholdConfig struct {
	Alert   int `yaml:"alert"`
	Warning int `yaml:"warning"`
	Error   int `yaml:"error"`
}
