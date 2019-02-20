package main
import "github.com/equalll/mydebug"

import (
	"flag"
	"fmt"
	"github.com/hidu/pproxy/serve"
	"log"
	"os"
)

var configPath = flag.String("conf", "./conf/pproxy.conf", "pproxy's config file")
var port = flag.Int("port", 0, "proxy port")
var vv = flag.Bool("vv", false, "debug,log request with more detail")
var showConf = flag.Bool("demo_conf", false, "show default conf")

var version = flag.Bool("v", false, "show version")

func init() {mydebug.INFO()
	df := flag.Usage

	flag.Usage = func() {
		df()
		fmt.Fprintln(os.Stderr, "\n HTTP protocol analysis tool\n https://github.com/hidu/pproxy/\n")
	}
}

func main() {mydebug.INFO()
	flag.Parse()

	if *showConf {
		demoConf := serve.GetDemoConf()
		fmt.Println(demoConf)
		os.Exit(0)
	}

	if *version {
		fmt.Println("pproxy version:", serve.GetVersion())
		os.Exit(0)
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Ldate)
	ser, err := serve.NewProxyServe(*configPath, *port)
	if err != nil {
		fmt.Println("start pproxy failed", err)
		os.Exit(2)
	}
	ser.Debug = *vv
	ser.Start()
}
