# process-runner
A powerful wrapper around os/exec `Cmd`. With support for channels.

Inspiration from https://github.com/jesselucas/executil

# Usage

# Simple example

```
outputBytes := command.Command("sh", "-c", "echo", "hallo").MustCombinedOutput()
fmt.Println("The output was:", string(outputBytes))
```

## Example with channels and prefix

```
outChannel := make(chan string)
errChannel := make(chan string)

cmd := command.Command("sh", "-c", "echo", "hallo")
cmd.OutputPrefix = "[EXAMPLE]"
cmd.StdoutChannel = outChannel
cmd.StderrChannel = errChannel

go func() {
    for {
        select {
        case outMsg := <-outChannel:
            fmt.Println("OUT: " + outMsg)
        case errMsg := <-errChannel:
            fmt.Println("ERR: " + errMsg)
        }
    }
}()

cmd.MustRun() // will panic if error occurred
```