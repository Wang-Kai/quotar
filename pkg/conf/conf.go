package conf

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

/*
	conf package will load all configuration, which will be used by service, including:

	1. workspace, which is nfs export path, and xfs mount pointer
	2. port, the port for gRPC server
*/

var WORKSPACE string
var PORT int

func init() {
	// set flag
	flag.String("workspace", "/data", "nfs export path, and xfs mount pointer")
	flag.Int("port", 10013, "the port for gRPC server")

	// bind flag to pflag
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// bind pflag to viper
	viper.BindPFlags(pflag.CommandLine)

	// catch env variables
	viper.BindEnv("PORT")
	viper.BindEnv("WORKSPACE")

	// set global variables
	WORKSPACE = viper.GetString("WORKSPACE")
	PORT = viper.GetInt("PORT")

	log.WithFields(log.Fields{
		"WORKSPACE": WORKSPACE,
		"PORT":      PORT,
	}).Info("Global variables catch by viper")
}
