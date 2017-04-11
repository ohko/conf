package conf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Conf 配置对象结构
type Conf struct {
	j map[string]interface{}
}

// NewConf 创建配置对象
func NewConf(confFile string) (*Conf, error) {

	o := &Conf{}
	o.j = make(map[string]interface{})

	if err := o.OpenConf(confFile); err != nil {
		o.j = make(map[string]interface{})
		return o, err
	}

	return o, nil
}

// OpenConf 打开配置文件
func (o *Conf) OpenConf(confFile string) error {
	if confFile == "" {
		return errors.New("confFile empty")
	}
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

// SetConf 设置配置内容
func (o *Conf) SetConf(data []byte) error {
	// 替换每行以 / 开头的内容为空
	// 去掉json中的注释
	reg, err := regexp.Compile("[ *|\\t*]// .*")
	if err != nil {
		return err
	}
	d := reg.ReplaceAll(data, nil)

	if err := json.Unmarshal(d, &o.j); err != nil {
		return err
	}
	return nil
}

// Clear 清空
func (o *Conf) Clear() {
	o.j = make(map[string]interface{})
}

// Exists 判断Key是否存在
func (o *Conf) Exists(key string) bool {
	if false == o.getSub(o.j, key, false) {
		return false
	}
	return true
}

// Get 获取
func (o *Conf) Get(key string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return o.getSub(o.j, key, defaultValue[0])
	}
	return o.getSub(o.j, key, nil)
}

// GetString 获取字符串值
func (o *Conf) GetString(key string, defaultValue ...string) string {
	if len(defaultValue) > 0 {
		return o.getSub(o.j, key, defaultValue[0]).(string)
	}
	return o.getSub(o.j, key, "").(string)
}

// GetInt 获取数字值
func (o *Conf) GetInt(key string, defaultValue ...int) int {
	var v interface{}
	if len(defaultValue) > 0 {
		v = o.getSub(o.j, key, defaultValue[0])
	} else {
		v = o.getSub(o.j, key, 0)
	}
	switch v.(type) {
	case float64:
		return int(v.(float64))
	default:
		return v.(int)
	}
}

// GetFloat64 获取浮点数
func (o *Conf) GetFloat64(key string, defaultValue ...float64) float64 {
	var v interface{}
	if len(defaultValue) > 0 {
		v = o.getSub(o.j, key, defaultValue[0])
	} else {
		v = o.getSub(o.j, key, 0)
	}
	switch v.(type) {
	case float64:
		return v.(float64)
	default:
		return float64(v.(int))
	}
}

// GetBool 获取是否
func (o *Conf) GetBool(key string, defaultValue ...bool) bool {
	if len(defaultValue) > 0 {
		return o.getSub(o.j, key, defaultValue[0]).(bool)
	}
	return o.getSub(o.j, key, false).(bool)
}

// Set 设置值
func (o *Conf) Set(key string, value interface{}) {
	o.setSub(o.j, key, value)
}

// ToBytes 生成bytes
func (o *Conf) ToBytes() []byte {
	bs, _ := json.Marshal(o.j)
	return bs
}

// ToString 生成字符串
func (o *Conf) ToString(prefix, indent string) string {
	bs, _ := json.MarshalIndent(o.j, prefix, indent)
	return string(bs)
}

// Map ...
func (o *Conf) Map() map[string]interface{} {
	return o.j
}

// Confs ...
func Confs(i interface{}, defaultValue ...[]*Conf) []*Conf {
	if i == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []*Conf{}
	}
	switch i.(type) {
	case []interface{}:
		r := []*Conf{}
		for _, v := range i.([]interface{}) {
			r = append(r, v.(*Conf))
		}
		return r
	default:
		return make([]*Conf, 0)
	}
}

// Interfaces ...
func Interfaces(i interface{}, defaultValue ...[]interface{}) []interface{} {
	if i == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []interface{}{}
	}
	switch i.(type) {
	case []interface{}:
		var r []interface{}
		for _, v := range i.([]interface{}) {
			r = append(r, v)
		}
		return r
	default:
		return make([]interface{}, 0)
	}
}

// Strings ...
func Strings(i interface{}, defaultValue ...[]string) []string {
	if i == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []string{}
	}
	switch i.(type) {
	case []interface{}:
		var r []string
		for _, v := range i.([]interface{}) {
			r = append(r, v.(string))
		}
		return r
	default:
		return make([]string, 0)
	}
}

// Ints ...
func Ints(i interface{}, defaultValue ...[]int) []int {
	if i == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []int{}
	}
	switch i.(type) {
	case []interface{}:
		var r []int
		for _, v := range i.([]interface{}) {
			r = append(r, int(v.(float64)))
		}
		return r
	default:
		return make([]int, 0)
	}
}

// Float64s ...
func Float64s(i interface{}, defaultValue ...[]float64) []float64 {
	if i == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []float64{}
	}
	switch i.(type) {
	case []interface{}:
		var r []float64
		for _, v := range i.([]interface{}) {
			r = append(r, v.(float64))
		}
		return r
	default:
		return make([]float64, 0)
	}
}

// Bools ...
func Bools(i interface{}, defaultValue ...[]bool) []bool {
	if i == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []bool{}
	}
	switch i.(type) {
	case []interface{}:
		var r []bool
		for _, v := range i.([]interface{}) {
			r = append(r, v.(bool))
		}
		return r
	default:
		return make([]bool, 0)
	}
}

func (o *Conf) setSub(obj map[string]interface{}, key string, value interface{}) map[string]interface{} {
	// 不包含"."
	if !strings.Contains(key, ".") {

		// 数组赋值
		if strings.Contains(key, "+") {
			_k := strings.Split(key, "+")
			_k1 := _k[0]
			_k2, _ := strconv.Atoi(_k[1])
			_new := false
			if _, ok := obj[_k1]; !ok {
				_new = true
				obj[_k1] = make([]interface{}, 1)
			}

			vv := obj[_k1].([]interface{})
			if _k[1] == "" && !_new { // 无序号，追加
				_k2 = len(vv)
			}
			// 填充
			if _k2+1 > len(vv) {
				_i := _k2 + 1 - len(vv)
				for i := 0; i < _i; i++ {
					vv = append(vv, nil)
				}
			}
			vv[_k2] = value
			obj[_k1] = vv

		} else { // 非数组赋值
			obj[key] = value
		}

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
		// 数组
		if strings.Contains(key, "+") {
			_k := strings.Split(key, "+")
			_k1 := _k[0]
			_k2, _ := strconv.Atoi(_k[1])

			// obj必定是map，如果有就返回值，没有返回默认值
			if v, ok := obj.(map[string]interface{})[_k1]; ok {
				vv := v.([]interface{})
				if _k2 < len(vv) {
					switch v.(type) {
					case map[string]interface{}:
						a, _ := NewConf("")
						a.j = vv[_k2].(map[string]interface{})
						return a
					default:
						return vv[_k2]
					}
				}
				return defaultValue
			}
		} else { // 非数组
			// obj必定是map，如果有就返回值，没有返回默认值
			if v, ok := obj.(map[string]interface{})[key]; ok {
				switch v.(type) {
				case []interface{}:
					r := []interface{}{}
					for _, vv := range v.([]interface{}) {
						switch vv.(type) {
						case map[string]interface{}:
							a, _ := NewConf("")
							a.j = vv.(map[string]interface{})
							r = append(r, a)
						default:
							r = append(r, vv)
						}
					}
					return r
				default:
					return v
				}
			}
			return defaultValue
		}
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
