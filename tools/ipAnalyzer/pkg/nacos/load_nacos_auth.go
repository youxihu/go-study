package nacos

import (
	"gopkg.in/yaml.v3"
	"ipAnalyzer/internal/entity"
	"os"
	"path/filepath"
)

// LoadNacosAuth 从本地文件加载 Nacos 认证配置
func LoadNacosAuth(filePath string) (*entity.AuthConfig, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	configFile, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var wrapper entity.AuthConfigWrapper
	err = yaml.Unmarshal(configFile, &wrapper)
	if err != nil {
		return nil, err
	}
	return &wrapper.Auth, nil
}
