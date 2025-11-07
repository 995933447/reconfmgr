### 是一个轻量级、可扩展的 Go 配置管理框架。
它通过统一的 Config 接口和中心化实现配置注册、通知驱动的热重载载，以及线程安全的内存存储。非常适合用于管理复杂业务或者大型复杂系统的配置。

### ✨ 特性

✅ 统一的配置接口 Config

🔁 支持任意配置源（文件、数据库、etcd、环境变量等）

🚦 支持通过触发配置热加载

🧱 模块化、可插拔的配置结构

📦 轻量依赖、无侵入集成

### 📦 安装
go get github.com/995933447/reconfmgr

### 🚀 快速上手
1️⃣ 定义你的配置实现

每个配置模块只需实现 Config 接口：
````
type Config interface {
	LoadConfig() error       // 首次加载配置
	ReloadConfig()           // 重载配置
	GetListenKeys() []string // 监听的keys来重载配置,只要触发Reload的keys其中一个命中(即交集)就会触发重载
	SetName(string)          // 设置配置唯一标识
	GetName() string         // 获取获赔唯一标识
	GetPriority() int        // 获取配置优先级，多个配置件套
}
````
可以嵌套ConfigBase来简化接口实现流程,ConfigBase做了基础的实现,具体组合的Config通常只需要实现GetListenKeys和LoadConfig方法即可。
或者具体的Config需要覆盖ConfigBase的默认行为。调用ConfigBase的ReloadConfig会默认调用组合的具体Config的Reload。

示例：
````
type PrintConfig struct {
	reconfmgr.ConfigBase
}

func (c *PrintConfig) GetPriority() int {
	return 1
}

func (c *PrintConfig) LoadConfig() error {
	fmt.Println("PrintConfig.LoadConfig", c.GetName())
	return nil
}

func (c *PrintConfig) GetListenKeys() []string {
	return []string{"tb" + c.GetName()}
}

func TestConfig(t *testing.T) {
	err := reconfmgr.Register("c1", &PrintConfig{})
	if err != nil {
		t.Errorf("register err:%v", err)
	}

	err = reconfmgr.Register("c2", &PrintConfig{})
	if err != nil {
		t.Errorf("register err:%v", err)
	}

	reconfmgr.Reload([]string{"tbc1", "tbc3"})
}
````
````
=== RUN   TestConfig
PrintConfig.LoadConfig c1
PrintConfig.LoadConfig c2
PrintConfig.LoadConfig c1
--- PASS: TestConfig (0.00s)
PASS
````
ConfigBase的实现代码：
````
var _ Config = (*ConfigBase)(nil)

type ConfigBase struct {
	name string
}

func (c *ConfigBase) GetPriority() int { // 默认间隔时间60s
	return 0
}

func (c *ConfigBase) SetName(name string) {
	c.name = name
}
func (c *ConfigBase) GetName() string {
	return c.name
}

func (c *ConfigBase) GetListenKeys() []string {
	return nil
}

func (c *ConfigBase) LoadConfig() (err error) {
	return nil
}

func (c *ConfigBase) ReloadConfig() {
	var cc Config
	if c.name != "" {
		var ok bool
		cc, ok = Get(c.name)
		if !ok {
			LogError("config " + c.name + " not found")
			return
		}
	}
	err := cc.LoadConfig()
	if err != nil {
		LogErrorf("config %s load config failed, err:%v", c.name, err)
	}
}
````

### 设计理念

reconfmgr = Simple, Unified, Reactive

reconfmgr 的目标不是提供具体配置格式解析,而是提供统一的注册与热更新机制让应用内部配置管理更加模块化和可维护。