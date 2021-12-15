package executors

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"wailik.com/internal/pkg/log"
)

type result struct {
	success bool
	output  string
}

type ShellExecutor struct{}

func (e *ShellExecutor) Execute(config map[string]interface{}, data interface{}, start Start, finish Finish) error {
	go func() {
		start()
		log.Debugf("execute command config(%+v) data(%+v)", config, data)
		r := executeCommand(buildCommand(config, data))
		finish(r.success, r.output, time.Now())
	}()

	return nil
}

func read(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser, resultChan chan result) {
	var outputStr string
	r := false
	scanner := bufio.NewScanner(std)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			resultChan <- result{
				success: r,
				output:  "timeout cancel",
			}

			return
		default:
			if scanner.Scan() {
				outputStr += scanner.Text()
			} else {
				err := scanner.Err()
				if err != nil {
					resultChan <- result{success: false, output: err.Error()}

					return
				}
				resultChan <- result{success: true, output: outputStr}

				return
			}
		}
	}
}

func executeCommand(cmd string) result {
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancelFunc context.CancelFunc) {
		time.Sleep(3 * time.Second)
		cancelFunc()
	}(cancel)

	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.CommandContext(ctx, "cmd", "/C", cmd) // windows
	} else {
		c = exec.CommandContext(ctx, "bash", "-c", cmd) // mac linux
	}

	stdout, err := c.StdoutPipe()
	if err != nil {
		return result{success: false, output: err.Error()}
	}

	stderr, err := c.StderrPipe()
	if err != nil {
		return result{success: false, output: err.Error()}
	}

	stdoutChan := make(chan result, 1)
	stderrChan := make(chan result, 1)
	defer func() {
		close(stderrChan)
		close(stdoutChan)
	}()

	var wg sync.WaitGroup

	wg.Add(2)
	go read(ctx, &wg, stderr, stderrChan)
	go read(ctx, &wg, stdout, stdoutChan)
	if err = c.Start(); err != nil {
		return result{success: false, output: err.Error()}
	}
	wg.Wait()

	if stderrReuslt := <-stderrChan; stderrReuslt.output != "" {
		return result{success: false, output: stderrReuslt.output}
	}

	return <-stdoutChan
}

func buildCommand(config map[string]interface{}, data interface{}) string {
	log.Info("cmd:" + config["cmd"].(string))
	log.Infof("data:%+v", data)

	return config["cmd"].(string)
}
