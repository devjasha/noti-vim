package main

import (
	"fmt"
	"os"

	"github.com/devjasha/noti-vim/internal/config"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	cfgFile string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "noti",
	Short:   "A fast CLI for managing Noti markdown notes",
	Long:    `Noti CLI is a standalone tool for managing markdown notes with Git integration`,
	Version: version,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/noti/config.yaml)")
	rootCmd.PersistentFlags().Bool("json", false, "output in JSON format")
	rootCmd.PersistentFlags().Bool("quiet", false, "minimal output")
}

func initConfig() {
	if err := config.Load(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
	}
}
