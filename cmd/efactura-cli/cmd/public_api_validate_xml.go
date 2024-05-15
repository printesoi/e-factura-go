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
	"context"
	"fmt"
	"os"

	"github.com/printesoi/e-factura-go/efactura"
	"github.com/spf13/cobra"
)

const (
	flagNameValidateInFilePath = "xml-file"
	flagNameValidateStandard   = "standard"
)

// publicApiValidateCmd represents the `public-api validate-xml` command
var publicApiValidateCmd = &cobra.Command{
	Use:   "validate-xml",
	Short: "Validate XML",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newEfacturaPublicClient(cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvXmlFile, err := cmd.Flags().GetString(flagNameValidateInFilePath)
		if err != nil {
			return err
		}
		f, err := os.Open(fvXmlFile)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		defer f.Close()

		fvStandard, err := cmd.Flags().GetString(flagNameValidateStandard)
		if err != nil {
			return err
		}
		switch fvStandard {
		case string(efactura.ValidateStandardFACT1), string(efactura.ValidateStandardFCN):
		default:
			return fmt.Errorf("invalid standard: `%s`", fvStandard)
		}

		validateRes, err := client.ValidateXML(context.Background(), f, efactura.ValidateStandard(fvStandard))
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		if !validateRes.IsOk() {
			cmd.SilenceUsage = true
			return fmt.Errorf("validate XML %s failed: %s", fvStandard, validateRes.GetFirstMessage())
		}
		fmt.Printf("validate %s: OK\n", fvStandard)
		return nil
	},
}

func init() {
	publicApiValidateCmd.Flags().StringP(flagNameValidateInFilePath, "f", "", "Path of the input XML file")
	_ = publicApiValidateCmd.MarkFlagRequired(flagNameValidateInFilePath)

	publicApiValidateCmd.Flags().StringP(flagNameValidateStandard, "s", "", fmt.Sprintf("Standard for validation (%s/%s)",
		efactura.ValidateStandardFACT1, efactura.ValidateStandardFCN))
	_ = publicApiValidateCmd.MarkFlagRequired(flagNameValidateStandard)

	publicApiCmd.AddCommand(publicApiValidateCmd)
}
