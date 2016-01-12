package command

import (
	"bufio"
	"io"
	"os/exec"
)

//Cmd is a wrapper around Cmd
type Cmd struct {
	StdoutChannel chan string
	StderrChannel chan string
	Cmd           *exec.Cmd
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
	// set the exec.Cmd
	c.Cmd = exec.Command(name, arg...)
	return c
}

// CombinedOutput wrapper for exec.Cmd.CombinedOutput()
func (c *Cmd) CombinedOutput() ([]byte, error) {
	return c.Cmd.CombinedOutput()
}
func (c *Cmd) MustCombinedOutput() []byte {
	out, err := c.CombinedOutput()
	checkError(err)
	return out
}

// Output wrapper for exec.Cmd.Output()
func (c *Cmd) Output() ([]byte, error) {
	return c.Cmd.Output()
}
func (c *Cmd) MustOutput() []byte {
	out, err := c.Output()
	checkError(err)
	return out
}

// Run wrapper for exec.Cmd.Run()
func (c *Cmd) Run() error {
	return c.Cmd.Run()
}
func (c *Cmd) MustRun() {
	checkError(c.Run())
}

// Start wrapper for exec.Cmd.Start()
func (c *Cmd) Start() error {
	// go routines to scan command out and err
	err := c.createPipeScanners()
	checkError(err)

	return c.Cmd.Start()
}
func (c *Cmd) MustStart() {
	checkError(c.Start())
}

// StderrPipe wrapper for exec.Cmd.StderrPipe()
func (c *Cmd) StderrPipe() (io.ReadCloser, error) {
	return c.Cmd.StderrPipe()
}
func (c *Cmd) MustStderrPipe() io.ReadCloser {
	pipe, err := c.StderrPipe()
	checkError(err)
	return pipe
}

// StdinPipe wrapper for exec.Cmd.StdinPipe()
func (c *Cmd) StdinPipe() (io.WriteCloser, error) {
	return c.Cmd.StdinPipe()
}
func (c *Cmd) MustStdinPipe() io.WriteCloser {
	pipe, err := c.StdinPipe()
	checkError(err)
	return pipe
}

// StdoutPipe wrapper for exec.Cmd.StdoutPipe()
func (c *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return c.Cmd.StdoutPipe()
}
func (c *Cmd) MustStdoutPipe() io.ReadCloser {
	pipe, err := c.StdoutPipe()
	checkError(err)
	return pipe
}

// Wait wrapper for exec.Cmd.Wait()
func (c *Cmd) Wait() error {
	return c.Cmd.Wait()
}
func (c *Cmd) MustWait() {
	checkError(c.Wait())
}

// Create stdout, and stderr pipes for given *Cmd
// Only works with cmd.Start()
func (c *Cmd) createPipeScanners() error {
	if c.StdoutChannel != nil {
		stdout, err := c.Cmd.StdoutPipe()
		if err != nil {
			return err
		}

		outScanner := bufio.NewScanner(stdout)

		go func() {
			for outScanner.Scan() {
				c.StdoutChannel <- outScanner.Text()
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
				c.StderrChannel <- errScanner.Text()
			}
		}()
	}

	return nil
}
