package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/goma"
	_ "github.com/cybozu-go/goma/actions/all"
	_ "github.com/cybozu-go/goma/filters/all"
	"github.com/cybozu-go/goma/monitor"
	_ "github.com/cybozu-go/goma/probes/all"
	"github.com/cybozu-go/log"
)

const (
	defaultConfDir    = "/usr/local/etc/goma"
	defaultListenAddr = "localhost:3838"
)

var (
	confDir    = flag.String("d", defaultConfDir, "directory for monitor configs")
	listenAddr = flag.String("s", defaultListenAddr, "HTTP server address")
)

func usage() {
	fmt.Fprint(os.Stderr, `Usage: goma [options] COMMAND [arg...]

If COMMAND is "serve", goma runs in server mode.
For other commands, goma works as a client for goma server.

Options:
`)
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, `
Commands:
    serve              Start agent server.
    list               List registered monitors.
    register FILE      Register monitors defined in FILE.
                       If FILE is "-", goma reads from stdin.
    show ID            Show the status of a monitor for ID.
    start ID           Start a monitor.
    stop ID            Stop a monitor.
    unregister ID      Stop and unregister a monitor.
    verbosity [LEVEL]  Query or change logging threshold.
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	cmd.LogConfig{}.Apply()

	args := flag.Args()

	if len(args) == 0 {
		usage()
		return
	}

	command := args[0]

	if command != "serve" {
		err := runCommand(command, args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, strings.TrimSpace(err.Error()))
			os.Exit(1)
		}
		return
	}

	if err := loadConfigs(*confDir); err != nil {
		log.ErrorExit(err)
	}

	goma.Serve(*listenAddr)
	err := cmd.Wait()
	if err != nil && !cmd.IsSignaled(err) {
		log.ErrorExit(err)
	}

	// stop all monitors gracefully.
	for _, m := range monitor.ListMonitors() {
		m.Stop()
	}
}
