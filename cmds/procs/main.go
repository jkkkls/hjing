package main

import (
	"flag"
	"log"

	"github.com/ochinchina/go-daemon"
)

var (
	// ./procs -c procs.yaml -d
	configName = flag.String("c", "procs.yaml", "config file name")
	daemonize  = flag.Bool("d", false, "run as daemon")
)

func main() {
	flag.Parse()
	conf, err := LoadConf(*configName)
	if err != nil {
		panic(err)
	}
	if *daemonize {
		Daemon("procs.log", func() error {
			return Run(conf)
		})
	} else {
		Run(conf)
	}
}

func Daemon(logfile string, f func() error) {
	ctx := &daemon.Context{LogFileName: logfile, PidFileName: "procs.pid"}
	d, err := ctx.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer ctx.Release()
	f()
}

func Run(conf *RunConfig) error {
	// for _, p := range conf.Procs {
	// 	// start
	// 	// stop
	// 	// restart
	// }
	return nil
}
