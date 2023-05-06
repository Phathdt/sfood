package common

import "github.com/namsral/flag"

var (
	ConsulHost = ""
)

func init() {
	flag.StringVar(&ConsulHost, "consul_host", "localhost:8500", "consult host, should be localhost:8500")
	flag.Parse()
}
