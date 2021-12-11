package executors

type Executor interface {
	Execute(config map[string]interface{}, data interface{}) error
}

func Get(name string) Executor {
	switch name {
	case "shell":
		return &ShellExecutor{}
	default:
		return nil
	}
}
