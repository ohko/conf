package conf

import "testing"

func TestConf(t *testing.T) {
	// 测试无配置文件
	cf, err := NewConf("")
	if err != nil {
		t.Log(err)
	}
	if cf.Exists("a.b.c") {
		t.Error(cf.Exists("a.b.c"))
	}
	if cf.GetInt("a", 1) != 1 {
		t.Error(1, cf.GetInt("a", 1))
	}
	if cf.GetString("b.c", "BC") != "BC" {
		t.Error(cf.GetString("b.c", "BC"))
	}

	// 测试有无内容
	cf.SetConf([]byte(`
{
	// 这是注释信息1
    "a": "A",
    "b": 2,
	// 这是注释信息2
    "c": { // 任意行中：空格+双斜线+空格，这个也是注释
        "d": 4.5,
        "e": {
            "f": 6.789,
			"g":"G",
			"gg":"// 双斜线前没有空格，这个应该保留",
			"h":true,
			"i":[1,2,3],
			"j":["a","b"],
			"k":[
				{
					"l":1,
					"m":2
				}
			]
        }
    }
}`))
	// a=A
	if cf.GetString("a") != "A" {
		t.Error(cf.GetString("a"))
	}
	// b=2
	if cf.GetInt("b") != 2 {
		t.Error(cf.GetInt("b"))
	}
	// c.d=4.5
	if cf.GetFloat64("c.d") != 4.5 {
		t.Error(cf.GetFloat64("c.d"))
	}
	// c.e.f=6.789
	if cf.GetFloat64("c.e.f") != 6.789 {
		t.Error(cf.GetFloat64("c.e.f"))
	}
	// c.e.g=G
	if cf.GetString("c.e.g") != "G" {
		t.Error(cf.GetString("c.e.g"))
	}
	// c.e.h=true
	if !cf.GetBool("c.e.h") {
		t.Error(cf.GetBool("c.e.h"))
	}
	// c.e.i=[1,2,3]
	if Ints(cf.Get("c.e.i"))[2] != 3 {
		t.Error(cf.Get("c.e.i"))
	}
	// c.e.j=["a","b"]
	if Strings(cf.Get("c.e.j"))[1] != "b" {
		t.Error(cf.Get("c.e.j"))
	}
	// c.e.k.m=true
	if Confs(cf.Get("c.e.k"))[0].GetInt("m") != 2 {
		t.Error(cf.Get("c.e.k"))
	}

	// 测试清空
	cf.Clear()
	if string(cf.ToBytes()) != `{}` {
		t.Error(string(cf.ToBytes()))
	}

	// 测试设置
	cf.Set("a", "AA")
	cf.Set("b.c", "BC")
	cf.Set("b.e.d", "BED")
	cf.Set("b.e.f", 123)
	cf.Set("b.e.g", 45.67)
	cf.Set("b.e.h+", 78.90)
	cf.Set("b.e.h+3", 78.93)
	cf.Set("b.e.h+4", 78.94)
	cf.Set("b.e.h+6", 78.96)
	cf.Set("b.e.h+", 78.97)
	cf1, _ := NewConf("")
	cf1.Set("a1", "A1")
	cf1.Set("a2", 2)
	cf.Set("b.e.j+", cf1.Map())
	if cf.GetString("a", "") != "AA" {
		t.Error(cf.GetString("a", ""))
	}
	if cf.GetString("b.c", "") != "BC" {
		t.Error(cf.GetString("b.c", ""))
	}
	if cf.GetInt("b.e.f", 0) != 123 {
		t.Error(cf.GetInt("b.e.f", 0))
	}
	if cf.GetFloat64("b.e.g", 0) != 45.67 {
		t.Error(cf.GetFloat64("b.e.g", 0))
	}
	v := cf.Get("b.e.j")
	if Confs(v)[0].GetString("a1", "") == "" {
		t.Error(Confs(v)[0].GetString("a1", ""))
	}
	if cf.GetFloat64("b.e.h+", 0) != 78.90 {
		t.Error(cf.GetFloat64("b.e.h+", 0))
	}
	if cf.GetFloat64("b.e.h+7", 0) != 78.97 {
		t.Error(cf.GetFloat64("b.e.h+7", 0))
	}
	if cf.GetFloat64("b.e.h+8", 0) != 0 {
		t.Error(cf.GetFloat64("b.e.h+8", 0))
	}

	// 结果判断
	if string(cf.ToBytes()) != `{"a":"AA","b":{"c":"BC","e":{"d":"BED","f":123,"g":45.67,"h":[78.9,null,null,78.93,78.94,null,78.96,78.97],"j":[{"a1":"A1","a2":2}]}}}` {
		t.Error(string(cf.ToBytes()))
	}
}
