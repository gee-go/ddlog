package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gee-go/ddlog/ddlog"
)

type Opts struct {
	ddlog.Config

	// How often to send logs.
	Rate time.Duration
}

func parseFlags() *Opts {
	o := &Opts{}
	flag.StringVar(&o.LogFormat, "fmt", ddlog.DefaultLogFormat, "a")
	flag.StringVar(&o.TimeFormat, "time", ddlog.DefaultTimeFormat, "a")
	var rate string
	flag.StringVar(&rate, "rate", "1s", "a")
	flag.Parse()

	r, err := time.ParseDuration(rate)
	if err != nil {
		log.Fatal(err)
	}
	o.Rate = r

	return o
}

func main() {
	opts := parseFlags()
	g := ddlog.NewGenerator(&opts.Config)
	g.UseUnicode = false

	for range time.Tick(opts.Rate) {
		fmt.Println(g.RandMsg())
	}
}
