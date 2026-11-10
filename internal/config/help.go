package config

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Printf(`
%s v[%s]

 -c, --config <file>  use this configuration file (built-in defaults apply first)
 -v, --debug          force debug level on
 -vv, --trace         force trace level on
 --nocolor            disable color and banner
 --banner             show banner + version and exit
 --genconfig          write default config to ./config.toml then exit
 -h, --help           show this help and exit
`, Title, Version)
	os.Exit(0)
}
