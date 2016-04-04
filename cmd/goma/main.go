package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cybozu-go/goma"
	_ "github.com/cybozu-go/goma/actions/all"
	_ "github.com/cybozu-go/goma/filters/all"
	"github.com/cybozu-go/goma/monitor"
	_ "github.com/cybozu-go/goma/probes/all"
	"github.com/cybozu-go/log"
	"golang.org/x/net/context"
)

const (
	defaultConfDir    = "/usr/local/etc/goma"
	defaultListenAddr = "localhost:3838"
)

var (
	confDir = flag.String("d", defaultConfDir, "directory for monitor configs")

	logLevel = flag.String("loglevel", "info", "logging level")
	logFile  = flag.String("logfile", "", "log filename")

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

	args := flag.Args()

	if len(args) == 0 {
		usage()
		return
	}

	cmd := args[0]

	if cmd != "serve" {
		err := runCommand(cmd, args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, strings.TrimSpace(err.Error()))
			os.Exit(1)
		}
		return
	}

	if err := log.DefaultLogger().SetThresholdByName(*logLevel); err != nil {
		log.ErrorExit(err)
	}

	if len(*logFile) > 0 {
		f, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.ErrorExit(err)
		}
		defer f.Close()
		log.DefaultLogger().SetOutput(f)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	if err := loadConfigs(ctx, *confDir); err != nil {
		log.ErrorExit(err)
	}

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.ErrorExit(err)
	}

	go func() {
		done <- goma.Serve(ctx, l)
	}()

	sig := make(chan os.Signal, 10)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	signal.Stop(sig)
	cancel()
	if err := <-done; err != nil {
		log.Error(err.Error(), nil)
	}

	// stop all monitors gracefully.
	for _, m := range monitor.ListMonitors() {
		m.Stop()
	}
}
