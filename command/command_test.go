package command

import (
	. "github.com/smartystreets/goconvey/convey"
	"runtime"
	"strings"
	"testing"
)

func createCommand(exe string, args ...string) *Cmd {
	switch runtime.GOOS {
	case "windows":
		allArgs := []string{}
		allArgs = append(allArgs, "/c")
		allArgs = append(allArgs, exe)
		allArgs = append(allArgs, args...)
		return Command("cmd", allArgs...)
	case "darwin", "linux":
		allArgs := []string{}
		allArgs = append(allArgs, "-c")
		allArgs = append(allArgs, exe)
		allArgs = append(allArgs, args...)
		return Command("sh", allArgs...)
	default:
		return Command(exe, args...)
	}
}

func sanitizeOutputString(s string) string {
	o := s
	o = strings.Replace(string(o), "\r", "", -1)
	o = strings.Replace(string(o), " \n", "\n", -1)
	return o
}

func TestChannels(t *testing.T) {
	Convey("No channels", t, func() {
		cmd := createCommand("echo", "hallo1", "&", "echo", "hallo2", "&", "echo", "hallo3", "1>&2")

		combinedOutputBytes := cmd.MustCombinedOutput()
		So(sanitizeOutputString(string(combinedOutputBytes)), ShouldEqual, sanitizeOutputString("hallo1\nhallo2\nhallo3\n"))
	})

	Convey("Unbuffered channels", t, func() {
		outChan := make(chan string)
		errChan := make(chan string)

		cmd := createCommand("echo", "hallo1", "&", "echo", "hallo2", "&", "echo", "hallo2", "1>&2")
		cmd.StdoutChannel = outChan
		cmd.StderrChannel = errChan

		cmd.MustStart()

		outputs := []string{}
		errors := []string{}
		go func() {
			for {
				select {
				case outMsg := <-outChan:
					outputs = append(outputs, outMsg)
				case errMsg := <-errChan:
					errors = append(errors, errMsg)
				}
			}
		}()

		cmd.MustWait()

		So(len(outputs), ShouldEqual, 2)
		So(len(errors), ShouldEqual, 1)
	})

	Convey("Buffered channels", t, func() {
		outChan := make(chan string, 20)
		errChan := make(chan string, 20)

		cmd := createCommand("echo", "hallo1", "&", "echo", "hallo2", "&", "echo", "hallo2", "1>&2")
		cmd.StdoutChannel = outChan
		cmd.StderrChannel = errChan

		cmd.MustStart()

		outputs := []string{}
		errors := []string{}
		go func() {
			for {
				select {
				case outMsg := <-outChan:
					outputs = append(outputs, outMsg)
				case errMsg := <-errChan:
					errors = append(errors, errMsg)
				}
			}
		}()

		cmd.MustWait()

		So(len(outputs), ShouldEqual, 2)
		So(len(errors), ShouldEqual, 1)
	})
}
