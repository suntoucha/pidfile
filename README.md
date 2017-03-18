# pidfile

Tiny library which implements basic work with PID File.

`import "github.com/awskii/pidfile"`

As easy as:
```
package main

import "github.com/awskii/pidfile"

func main() {
    pf := pidfile.Init("name_of_pid_file")

    // Checking pid file. If there is no process with same PID,
    // err will be nil
    if err := pf.TryLock(); err != nil {
        fmt.Println(err)
        return
    }
    defer pf.Remove() // removing pid file after usage

    // some code here
}
```
