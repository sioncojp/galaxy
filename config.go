package galaxy

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config ...Load config.yml
type Config struct {
	WorkDir string
	Server  struct {
		Port int
	}
	Database struct {
		Driver   string
		User     string
		Password string
		Host     string
		Port     int
		DBname   string
	}
	Github struct {
		Repository string
		Name       string
	}
	Url    string
	Script string
	Docker struct {
		Image      string
		Tag        string
		ProxyImage string
		ProxyTag   string
		Exec       string
	}
}

// LoadConfig ...Load config.yml
func LoadConfig(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err = yaml.Unmarshal(buf, &c); err != nil {
		return nil, err
	}
	return c, nil
}
