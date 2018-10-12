package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/naoina/toml"

	"github.com/spf13/cobra"
	"github.com/xtimeline/gox/log"
)

var configFile string

var opts = &Options{
	LogLevel: "warn",
}

var RootCmd = &cobra.Command{
	Use:  "ztasker",
	Long: "ztasker",
}

type Options struct {
	LogLevel string `flag:"log-level" cfg:"log_level"`
}

func Execute(task Tasker) {
	runCmd.Run = func(cmd *cobra.Command, args []string) {
		if task != nil {
			if task.TaskConf != nil && len(args) != 0 {
				if err := parseConfig(args[0], task.TaskConf()); err != nil {
					panic(err)
				}
			}
			run(task, args)
		} else {
			os.Exit(-1)
		}
	}
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// RootCmd.PersistentFlags().StringVar(&configFile, "config", "/etc/cooper.toml", "path to config file")
}

func initConfig() {
	l.SetLevel(opts.LogLevel)
	l.SetOutput(os.Stderr)
}

func parseConfig(file string, injected interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if err := toml.Unmarshal(buf, injected); err != nil {
		return err
	}
	return nil
}
