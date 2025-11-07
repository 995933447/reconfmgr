package reconfmgr

import (
	"fmt"
	"sort"
)

func isListenKey(c Config, key string) bool {
	for _, v := range c.GetListenKeys() {
		if v == "*" || v == key {
			return true
		}
	}
	return false
}

func Reload(keys []string) {
	LogInfof("============================Begin ReloadConfigs for keys:%+v============================", keys)

	configMu.Lock()

	var (
		shouldReloadConfigs   []Config
		shouldReloadConfigSet = map[Config]struct{}{}
	)
	for _, config := range configs {
		for _, key := range keys {
			if !isListenKey(config, key) {
				continue
			}

			shouldReloadConfigSet[config] = struct{}{}
			shouldReloadConfigs = append(shouldReloadConfigs, config)
			break
		}
	}

	configMu.Unlock()

	if len(shouldReloadConfigs) == 0 {
		return
	}

	sort.Slice(shouldReloadConfigs, func(i, j int) bool {
		return shouldReloadConfigs[i].GetPriority() > shouldReloadConfigs[j].GetPriority()
	})

	for name, config := range shouldReloadConfigs {
		LogInfof("==============ReloadConfigCache name:%s priority:%d  listen:%+v==============", name, config.GetPriority(), config.GetListenKeys())
		config.ReloadConfig()
		LogInfof("==============end ReloadConfigCache [%s]  listen:%+v==============", name, config.GetListenKeys())
	}

	LogInfof(fmt.Sprintf("============================End ReloadConfigCache for keys:%+v============================", keys))
}
