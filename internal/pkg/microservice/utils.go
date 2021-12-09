package microservice

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"gopkg.in/yaml.v2"

	"wailik.com/internal/pkg/log"
)

func GetNodeNameByNodeId(id string) string {
	sli := strings.Split(id, "/")
	log.Debugf("%+v", len(sli))

	nodeName := sli[len(sli)-2]
	log.Debugf("%+v", nodeName)

	return nodeName
}

func LoadClientConfig(path string) (*ClientConfig, error) {
	conf := clientv3.Config{}
	yamlConf := make(map[interface{}]interface{})
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(f, &yamlConf); err != nil {
		return nil, err
	}

	log.Debugf("confY:%+v", yamlConf)

	jsonConf := make(map[string]interface{})
	for key, value := range yamlConf {
		switch key := key.(type) {
		case string:
			jsonConf[key] = value
		}
	}

	j, err := json.Marshal(&jsonConf)
	if err != nil {
		return nil, err
	}

	log.Debugf("json:%+v", string(j))

	if err := json.Unmarshal(j, &conf); err != nil {
		return nil, err
	}

	log.Debugf("config:%+v", conf)

	return (*ClientConfig)(&conf), nil
}
