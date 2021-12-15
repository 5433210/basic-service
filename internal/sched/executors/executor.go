package executors

import "time"

type Start func()

type Finish func(success bool, output string, ts time.Time)

type Executor interface {
	Execute(config map[string]interface{}, data interface{}, start Start, finish Finish) error
}

func Get(name string) Executor {
	switch name {
	case "shell":
		return &ShellExecutor{}
	default:
		return nil
	}
}
