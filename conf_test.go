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
    "a__comment": "A",
    "a": "A",
    "b__comment": 2,
    "b": 2,
    "c__comment": "C",
    "c": {
        "d__comment": 4.5,
        "d": 4.5,
        "e": {
            "f__comment": 6.789,
            "f": 6.789,
			"g":"G"
        }
    }
}`))
	// a=A
	if cf.GetString("a", "A") != "A" {
		t.Error("A", cf.GetString("a", "A"))
	}
	// b=2
	if cf.GetInt("b", 2) != 2 {
		t.Error(2, cf.GetInt("b", 2))
	}
	// c.d=4.5
	if cf.GetFloat64("c.d", 1.0) != 4.5 || cf.GetFloat64("c.d", 1.0) == 1.0 {
		t.Error(cf.GetFloat64("c.d", 1.0))
	}
	// c.e.f=6.789
	if cf.GetFloat64("c.e.f", 1.0) != 6.789 || cf.GetFloat64("c.e.f", 1.0) == 1.0 {
		t.Error(cf.GetFloat64("c.e.f", 1.0))
	}
	// c.e.g=G
	if cf.GetString("c.e.g", "A") != "G" || cf.GetString("c.e.g", "A") == "A" {
		t.Error(cf.GetString("c.e.g", "A"))
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

	// 结果判断
	if string(cf.ToBytes()) != `{"a":"AA","b":{"c":"BC","e":{"d":"BED","f":123,"g":45.67}}}` {
		t.Error(string(cf.ToBytes()))
	}
}
