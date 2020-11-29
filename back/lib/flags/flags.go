// This package assists in reading command line flags.
package flags

import (
	"flag"

	"github.com/FTi130/keep-the-moment-app/back/lib/config"
	"github.com/FTi130/keep-the-moment-app/back/server"
)

// Flags structure contains variables, which should be received as commandline flags.
type Flags struct {
	Config config.Flags
	Server server.Flags
}

// Read reads commandline flags.
func Read() (f Flags) {
	f.Config.Path = flag.String("config", "./", "path with config.toml")
	f.Server.Debug = flag.Bool("debug", false, "debug mode")

	flag.Parse()
	return
}
