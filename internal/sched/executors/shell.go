package executors

import (
	"time"

	"wailik.com/internal/pkg/log"
)

type ShellExecutor struct{}

func (e *ShellExecutor) Execute(config map[string]interface{}, data interface{}, start Start, finish Finish) error {
	go func() {
		start()
		log.Debugf("execute config(%+v) data(%+v)", config, data)
		finish("0", "ok", time.Now())
	}()

	return nil
}
