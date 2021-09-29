package options

import (
	"flag"
	"log"
	"os"
)

type Options struct {
	Port        int
	Password    string
	Servername  string
	NetworkFile string
	Verbosity   bool
}

func (opts *Options) GetOptions(log *log.Logger) {
	flag.IntVar(&opts.Port, "p", 6697, "port number")
	flag.StringVar(&opts.Password, "password", "", "server password")
	flag.StringVar(&opts.Servername, "servname", "chirc", "server name")
	flag.StringVar(&opts.NetworkFile, "networkfile", "", "network-file path")
	flag.BoolVar(&opts.Verbosity, "v", false, "verbose output")

	flag.Parse()

	// if opts.Password == "" {
	// 	fmt.Fprint(os.Stderr, "ERROR: password cant be empty (-password option)\n")
	// 	os.Exit(-1)
	// }

	if (opts.NetworkFile != "") && !(opts.Servername == "") {
		log.Println("ERROR: If specifying a network file, you must also specify a server name.")
		os.Exit(-1)
	}
}
