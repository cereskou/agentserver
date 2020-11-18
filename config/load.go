package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kardianos/osext"
)

//CONFIG -
var CONFIG *AppSetting

//GetAppSetting -
func GetAppSetting() *AppSetting {
	return CONFIG
}

//GetConfig -
func GetConfig(stage string) (*AppSetting, error) {
	if CONFIG == nil {
		exename, _ := osext.Executable()
		dir := filepath.Dir(exename)
		filename := filepath.Join(dir, "app.config")
		var conf AppSetting
		if _, err := toml.DecodeFile(filename, &conf); err != nil {
			return nil, err
		}

		CONFIG = &conf
	}

	return CONFIG, nil
}
