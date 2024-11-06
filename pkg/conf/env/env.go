package env

import (
	"flag"
	"os"
)

var (
	Hostname string // hostname
	Mode     string // run mode
)

const (
	DEBUG   = "debug"
	TEST    = "test"
	GRAY    = "gray"
	RELEASE = "release"
)

func init() {
	var err error
	if Hostname, err = os.Hostname(); err != nil || Hostname == "" {
		Hostname = os.Getenv("HOSTNAME")
	}
	addFlag(flag.CommandLine)
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}

// default < goenv < command
func addFlag(f *flag.FlagSet) {
	f.StringVar(&Mode, "mode", defaultString("MODE", DEBUG), "run mode,default debug mode.value:[debug,test,gray,release]")
}
