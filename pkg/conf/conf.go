package conf

import (
	"flag"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
	conf package will load all configuration, which will be used by service, including:

	1. workspace, which is nfs export path, and xfs mount pointer
*/

var (
	WORKSPACE string
)

func init() {
	// set flag
	flag.String("workspace", "/data", "nfs export path, and xfs mount pointer")

	// bind flag to pflag
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// bind pflag to viper
	viper.BindPFlags(pflag.CommandLine)

	// catch env variables
	viper.AutomaticEnv()
	if err := viper.BindEnv("WORKSPACE"); err != nil {
		log.Fatal(err)
	}

	// set global variables
	WORKSPACE = viper.GetString("WORKSPACE")
}
