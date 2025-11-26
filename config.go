package reconfmgr

import (
	"sync"
)

type Config interface {
	LoadConfig() error       // 首次加载配置
	ReloadConfig()           // 重载配置
	GetListenKeys() []string // 监听的keys来重载配置,只要触发Reload的keys其中一个命中(即交集)就会触发重载
	SetName(string)          // 设置配置唯一标识
	GetName() string         // 获取获赔唯一标识
	GetPriority() int        // 获取配置优先级，多个配置件套
}

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

var (
	configMu sync.RWMutex
	configs  = make(map[string]Config)
)

func Register(name string, config Config) error {
	config.SetName(name)
	err := config.LoadConfig()
	if err != nil {
		return err
	}
	configMu.Lock()
	defer configMu.Unlock()
	configs[name] = config
	return nil
}

func Get(name string) (Config, bool) {
	configMu.RLock()
	defer configMu.RUnlock()

	config, ok := configs[name]
	return config, ok
}

func MustGet(name string) Config {
	config, ok := Get(name)
	if !ok {
		panic("MustGetConfig: config not found: " + name)
	}
	return config
}
