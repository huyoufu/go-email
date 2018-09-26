package email

import (
	"github.com/go-yaml/yamL"
	"io/ioutil"
	"net/smtp"
	"os"
	"path/filepath"
)

type Config struct {
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	From     string `yaml:"from"`
}

func (c *Config) auth() smtp.Auth {
	return smtp.PlainAuth("", c.Account, c.Password, c.Host)
}
func parseYaml() (c *Config) {
	dir := filepath.Dir(os.Args[0])
	bytes, _ := ioutil.ReadFile(dir + "/email.yaml")
	c = &Config{}
	yaml.Unmarshal(bytes, c)
	return
}
