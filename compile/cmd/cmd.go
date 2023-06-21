package cmd

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	mutex    sync.RWMutex
	cmdStore = make(map[string]map[Images]Command)
)

func init() {
	Register(&ethSolidityCmdV060{})
	Register(&bcosSolidityCmdV060{})
	Register(&chainmakerSolidityCmdV060{})
	Register(&xuperSolidityCmdV060{})
}

func Register(o Command) error {
	name := o.Organization()
	if name == "" {
		return errors.New("organization is empty")
	}

	image := o.Images()
	if image == "" {
		return errors.New("images is empty")
	}

	mutex.Lock()
	defer mutex.Unlock()

	v, ok := cmdStore[name]
	if !ok {
		cmdStore[name] = make(map[Images]Command)
	}
	if _, ok := v[Images(image)]; ok {
		return fmt.Errorf("%s images is register", image)
	}
	cmdStore[name] = map[Images]Command{Images(image): o}
	return nil
}

func Get(org string) (map[Images]Command, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	v, ok := cmdStore[org]
	return v, ok
}

type Command interface {
	Organization() string
	Images() string
	Cmd(Args) (string, []string, error)
}

type Args struct {
	Organization string
	Images       Images
	HostDir      string
	TargetDir    string
	Overwrite    bool
	Extend       map[string]interface{}
}

func (c *Args) Valid() error {
	if c.Organization == "" {
		return errors.New("organization is empty")
	}
	if c.Images == "" {
		return errors.New("image is empty")
	}
	if c.HostDir == "" {
		return errors.New("hostDir is empty")
	}
	if c.TargetDir == "" {
		return errors.New("targetDir is empty")
	}
	return nil
}

//func (c *Args) Select(name Images) {
//	mutex.Lock()
//	defer mutex.Unlock()
//	v, ok := cmdStore[name]
//	if !ok {
//		return
//	}
//}

type Images string

func (i Images) String() string {
	return string(i)
}

func (i Images) Repository() string {
	v := strings.Split(i.String(), ":")
	switch len(v) {
	case 1:
		return i.String()
	case 2:
		return v[0]
	default:
		return ""
	}
}

func (i Images) Tag() string {
	v := strings.Split(i.String(), ":")
	switch len(v) {
	case 1:
		return "latest"
	case 2:
		return v[1]
	default:
		return ""
	}
}
