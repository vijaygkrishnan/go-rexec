# go-rexec
Basic script to execute CLI commands on multiple cisco nexus switches or linux hosts. Converted my python scripts to go to tryout go-language. Since go encourages and makes it easy to open source, here is my first github commit !

# build
Call Main() from go-rexec package and build to get a runnable binary

example/remote_exec.go:
package main

import rexec "github.com/vijaygkrishnan/go-rexec"

func main() {
            rexec.Main()
}

go build remote_exec.go

# usage
./remote_exec ?cmd ?show version | grep N9K?
./remote_exec ?sort ?passwd 12345 ?cmd  ?copy running-cfg startup-cfg" 
