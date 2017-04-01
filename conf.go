package conf

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Conf ...
type Conf struct {
	j map[string]interface{}
}

// NewConf ...
func NewConf(confFile string) (*Conf, error) {

	o := &Conf{}
	o.j = make(map[string]interface{})

	if err := o.OpenConf(confFile); err != nil {
		o.j = make(map[string]interface{})
		return o, err
	}

	return o, nil
}

// OpenConf ...
func (o *Conf) OpenConf(confFile string) error {
	// 读取配置文件
	confContent, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}

	// 解析配置
	if err := o.SetConf(confContent); err != nil {
		return err
	}

	return nil
}

// SetConf ...
func (o *Conf) SetConf(data []byte) error {
	if err := json.Unmarshal(data, &o.j); err != nil {
		return err
	}
	return nil
}

// Clear ...
func (o *Conf) Clear() {
	o.j = make(map[string]interface{})
}

// Get ...
func (o *Conf) Get(key string, defaultValue interface{}) interface{} {
	return o.getSub(o.j, key, defaultValue)
}

// GetString ...
func (o *Conf) GetString(key string, defaultValue string) string {
	return o.getSub(o.j, key, defaultValue).(string)
}

// GetInt ...
func (o *Conf) GetInt(key string, defaultValue int) int {
	v := o.getSub(o.j, key, defaultValue)
	switch v.(type) {
	case float64:
		return int(v.(float64))
	default:
		return v.(int)
	}
}

// GetFloat64 ...
func (o *Conf) GetFloat64(key string, defaultValue float64) float64 {
	v := o.getSub(o.j, key, defaultValue)
	switch v.(type) {
	case float64:
		return v.(float64)
	default:
		return float64(v.(int))
	}
}

// Set ...
func (o *Conf) Set(key string, value interface{}) {
	o.setSub(o.j, key, value)
}

// ToBytes ...
func (o *Conf) ToBytes() []byte {
	bs, _ := json.Marshal(o.j)
	return bs
}

// ToString ...
func (o *Conf) ToString(prefix, indent string) string {
	bs, _ := json.MarshalIndent(o.j, prefix, indent)
	return string(bs)
}

func (o *Conf) setSub(obj map[string]interface{}, key string, value interface{}) map[string]interface{} {
	// 不包含"."
	if !strings.Contains(key, ".") {
		obj[key] = value
		return obj
	}

	// 包含 ".": "a.b.c"
	_k := strings.Split(key, ".")
	_k1 := _k[0]                     // "a"
	_k2 := strings.Join(_k[1:], ".") // "b.c"
	// map中存在，就继续查找sub
	var sub map[string]interface{}
	if _, ok := obj[_k1]; !ok {
		sub = make(map[string]interface{})
		obj[_k1] = sub
	} else {
		sub = obj[_k1].(map[string]interface{})
	}

	return o.setSub(sub, _k2, value)
}

func (o *Conf) getSub(obj interface{}, key string, defaultValue interface{}) interface{} {
	// 不包含"."
	if !strings.Contains(key, ".") {
		// obj必定是map，如果有就返回值，没有返回默认值
		if v, ok := obj.(map[string]interface{})[key]; ok {
			return v
		}
		return defaultValue
	}

	// 包含 ".": "a.b.c"
	_k := strings.Split(key, ".")
	_k1 := _k[0]                     // "a"
	_k2 := strings.Join(_k[1:], ".") // "b.c"
	// map中存在，就继续查找sub
	if v, ok := obj.(map[string]interface{})[_k1]; ok {
		switch v.(type) {
		case map[string]interface{}: // 如果是map
			return o.getSub(v, _k2, defaultValue)
		default:
			return v
		}
	}

	return defaultValue
}
