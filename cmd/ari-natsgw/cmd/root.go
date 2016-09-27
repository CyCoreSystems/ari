package cmd

import (
	"fmt"
	"os"
	"strings"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/nats-io/nats"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ari-natsgw",
	Short: "A gateway that exposes ARI over nats",
	Long: `ari-natsgw is a daemon that connects to the NATS server and to
Asterisk via ARI. ARI commands are exposed over NATS for operation
via the ari client transport under github.com/CyCoreSystems/ari/client/nc.`,
	Run: func(cmd *cobra.Command, args []string) {

		// setup logging
		log := log15.New()
		var handler log15.Handler = log15.StdoutHandler

		if verbose {
			handler = log15.LvlFilterHandler(log15.LvlDebug, handler)
		} else {
			handler = log15.LvlFilterHandler(log15.LvlInfo, handler)
		}
		log.SetHandler(handler)

		// run server
		os.Exit(runServer(log))
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	p := RootCmd.PersistentFlags()

	p.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ari-natsgw.yaml)")
	p.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	p.String("nats.url", nats.DefaultURL, "URL for connecting to the NATS cluster")
	p.String("ari.application", "", "ARI Stasis Application")
	p.String("ari.username", "", "Username for connecting to ARI")
	p.String("ari.password", "", "Password for connecting to ARI")
	p.String("ari.http_url", "http://localhost:8088/ari", "HTTP Base URL for connecting to ARI")
	p.String("ari.websocket_url", "ws://localhost:8088/ari/events", "Websocket URL for connecting to ARI")

	for _, n := range []string{"verbose", "nats.url", "ari.application", "ari.username", "ari.password", "ari.http_url", "ari.websocket_url"} {
		viper.BindPFlag(n, p.Lookup(n))
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ari-natsgw") // name of config file (without extension)
	viper.AddConfigPath("$HOME")       // adding home directory as first search path
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
