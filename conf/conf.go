package conf

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Conf struct {
	OpenAI OpenAI `yaml:"openai"`
}

type OpenAI struct {
	Token string `yaml:"token"`
}

func New(conf *string) (*Conf, error) {
	c := Conf{
		OpenAI: OpenAI{
			Token: "",
		},
	}

	data, err := getYamlFile(*conf)
	if err != nil {
		return &c, err
	}

	y := []byte(os.ExpandEnv(string(data)))
	if err := yaml.Unmarshal(y, &c); err != nil {
		return &c, err
	}

	if c.OpenAI.Token == "" {
		return &c, fmt.Errorf("OpenAI token is empty")
	}

	log.Println("Load the configuration file successfully!")

	return &c, nil
}

func getYamlFile(path string) ([]byte, error) {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}

	if f.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}
	return os.ReadFile(path)
}
