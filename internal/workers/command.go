package workers

import (
	"fmt"
	"strings"
)

// CommandJob runs a command and wraps the output in a job
func (c *Client) CommandJob(binary string, args ...string) Job {
	list := append([]string{binary}, args...)
	return func() error {
		status, err := Shell(strings.Join(list, " "), ShellOptions{Out: c.ProcessLine, Err: c.ProcessError})
		if err != nil {
			return err
		}
		if status.Exit != 0 {
			return fmt.Errorf("command failed with exit code %d", status.Exit)
		}
		return nil
	}
}

func (c *Client) ProcessLine(name, line string) {
	// Do something with the line here.
	c.logger.Warnf("[%s] %s", name, line)
}

func (c *Client) ProcessError(name, line string) {
	// Do something with the line here.
	c.logger.Errorf("[%s] %s", name, line)
}

// https://www.dolthub.com/blog/2022-11-28-go-os-exec-patterns/
/*
	cmd := exec.Command("ls", "/usr/local/bin")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	if err != nil {
		return err
	}
	for scanner.Scan() {
		// Do something with the line here.
		ProcessLine(scanner.Text())
	}
	if scanner.Err() != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return scanner.Err()
	}
	return cmd.Wait()
*/
