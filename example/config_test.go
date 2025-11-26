package example

import (
	"fmt"
	"testing"

	"github.com/995933447/reconfmgr"
)

var _ reconfmgr.Config = (*reconfmgr.ConfigBase)(nil)

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
