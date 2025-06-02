// Copyright 2024 Victor Dodon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "efactura-cli",
	Short: "A CLI client for the ANAF e-factura APIs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

const (
	flagNameProduction = "production"
)

var (
	efacturaCfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&efacturaCfgFile, "config", "", "config file (default is $HOME/.config/e-factura.yaml)")
	rootCmd.PersistentFlags().Bool(flagNameProduction, false, "Production mode (default sandbox)")

	bindViperFlag := func(name string) {
		viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
		viper.BindEnv(name)
	}
	viper.SetEnvPrefix("EFACTURA")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	for _, flagName := range []string{
		flagNameProduction,
	} {
		bindViperFlag(flagName)
	}
}

func initConfig() {
	if efacturaCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(efacturaCfgFile)
	} else {
		// Find user config directory (on Unix systems - $XDG_CONFIG_DIR)
		config, err := os.UserConfigDir()
		cobra.CheckErr(err)

		// Search config in user config directory with name "e-factura.yaml".
		viper.AddConfigPath(config)
		viper.SetConfigType("yaml")
		viper.SetConfigName("e-factura")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// This is a hack to make cobra required flags work with viper
	// https://github.com/spf13/viper/issues/397#issuecomment-544272457
	postInitCommands(rootCmd.Commands())
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
