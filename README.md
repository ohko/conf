# Conf json格式的配置文件读写库，用于golang。

- 支持JSON格式
- 支持子节点获取
- 支持重置配置文件
- 支持清空配置
- 支持设置配置项
- 支持输出配置内容

# 开始使用

```
$ go get -u -v github.com/ohko/conf
```

# 打开配置文件

```
import "github.com/ohko/conf"
cf, _ := NewConf("conf.json")
```

# 判断配置项

```
if cf.Exists("a.b.c") {
    ...
}
```

# 获取配置项

```
i:=cf.GetInt("a", 0)
s:=cf.GetString("a.b", "")
```

# 设置配置内容

```
cf.SetConf([]byte(`{"a":"AA","b":{"c":"BC","e":{"d":"BED","f":123,"g":45.67}}}`))
```

# 设置配置项

```
cf.Set("a", "AA")
cf.Set("b.c", "BC")
cf.Set("b.e.d", "BED")
cf.Set("b.e.f", 123)
cf.Set("b.e.g", 45.67)
```

# 清空配置

```
cf.Clear()
```

# 输出配置

```
cf.ToBytes()
cf.ToString()
```