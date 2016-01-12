package command

import (
	"bufio"
	"io"
	"os/exec"
)

//Cmd is a wrapper around Cmd
type Cmd struct {
	*exec.Cmd
	createdPipes bool

	OutputPrefix  string
	StdoutChannel chan string
	StderrChannel chan string
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Command returns a an executil Cmd struct with
// an exec.Cmd struct embedded in it
func Command(name string, arg ...string) *Cmd {
	c := new(Cmd)
	c.OutputPrefix = ""
	c.Cmd = exec.Command(name, arg...)
	return c
}

// MustCombinedOutput calls CombinedOutput() and panics on error
func (c *Cmd) MustCombinedOutput() []byte {
	out, err := c.CombinedOutput()
	checkError(err)
	return out
}

// MustOutput calls Output() and panics on error
func (c *Cmd) MustOutput() []byte {
	out, err := c.Output()
	checkError(err)
	return out
}

// Run wrapper for exec.Cmd.Run() - which also calls `createPipeScanners` beforehand to setup pipes/channels
func (c *Cmd) Run() error {
	err := c.createPipeScanners()
	if err != nil {
		return err
	}
	return c.Cmd.Run()
}
func (c *Cmd) MustRun() {
	checkError(c.Run())
}

// Start wrapper for exec.Cmd.Start() - which also calls `createPipeScanners` beforehand to setup pipes/channels
func (c *Cmd) Start() error {
	// go routines to scan command out and err
	err := c.createPipeScanners()
	if err != nil {
		return err
	}

	return c.Cmd.Start()
}
func (c *Cmd) MustStart() {
	checkError(c.Start())
}

// MustStderrPipe calls StderrPipe() and panics on error
func (c *Cmd) MustStderrPipe() io.ReadCloser {
	pipe, err := c.StderrPipe()
	checkError(err)
	return pipe
}

// MustStdinPipe calls StdinPipe() and panics on error
func (c *Cmd) MustStdinPipe() io.WriteCloser {
	pipe, err := c.StdinPipe()
	checkError(err)
	return pipe
}

// MustStdoutPipe calls StdoutPipe() and panics on error
func (c *Cmd) MustStdoutPipe() io.ReadCloser {
	pipe, err := c.StdoutPipe()
	checkError(err)
	return pipe
}

// MustWait calls Wait() and panics on error
func (c *Cmd) MustWait() {
	checkError(c.Wait())
}

// Create stdout, and stderr pipes for given *Cmd
// Only works with cmd.Start()
func (c *Cmd) createPipeScanners() error {
	if c.createdPipes {
		return nil
	}

	if c.StdoutChannel != nil {
		stdout, err := c.Cmd.StdoutPipe()
		if err != nil {
			return err
		}

		outScanner := bufio.NewScanner(stdout)

		go func() {
			for outScanner.Scan() {
				c.StdoutChannel <- c.OutputPrefix + outScanner.Text()
			}
		}()
	}

	if c.StderrChannel != nil {
		stderr, err := c.Cmd.StderrPipe()
		if err != nil {
			return err
		}

		errScanner := bufio.NewScanner(stderr)

		// Scan for text
		go func() {
			for errScanner.Scan() {
				c.StderrChannel <- c.OutputPrefix + errScanner.Text()
			}
		}()
	}

	return nil
}
