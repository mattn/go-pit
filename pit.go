package pit

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type pit struct {
	directory string
	config    string
	profile   string
}

var instance *pit

func GetInstance() *pit {
	if instance == nil {
		d := path.Join(os.Getenv("HOME"), ".pit")
		instance = &pit{
			directory: d,
		}
		instance.SetProfile("default")
		instance.config = path.Join(d, "pit.yaml")
	}
	return instance
}

func (pit *pit) SetProfile(name string) {
	pit.profile = path.Join(pit.directory, name+".yaml")
}

func (pit pit) CurrentProfile() (profile string) {
	self := GetInstance()
	m := self.Config()
	profile = m["profile"].(string)
	return
}

func (pit pit) Load() (profile map[interface{}]interface{}) {
	pit.SetProfile(pit.CurrentProfile())

	// TODO: ファイル無いとき
	b, err := ioutil.ReadFile(pit.profile)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &profile)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (pit pit) Config() (profile map[interface{}]interface{}) {
	b, err := ioutil.ReadFile(pit.config)

	// TODO: ファイル無いとき
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &profile)
	if err != nil {
		log.Fatal(err)
	}
	return
}

type Profile map[string]string
type Config struct {
	Profile string `yaml:profile`
}

func Get(name string) (profile Profile) {
	self := GetInstance()
	m := self.Load()

	// これもちっとマシに型変換できないのかな...
	profile = make(Profile)
	for k, v := range m[name].(map[interface{}]interface{}) {
		profile[k.(string)] = v.(string)
	}
	return
}

func Switch(name string) (prev string) {
	self := GetInstance()
	self.SetProfile(name)

	config := self.Config()
	prev = config["profile"].(string)
	c := Config{
		Profile: name,
	}

	// FIXME: エラー無視してる
	b, _ := yaml.Marshal(&c)
	err := ioutil.WriteFile(self.config, b, 0600)
	if err != nil {
		log.Fatal(err)
	}
	return
}