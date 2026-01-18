// Copyright 2024-2026 Victor Dodon
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

	"github.com/spf13/cobra"
)

const (
	flagNameValidateSignatureInvoiceXML = "invoice-xml"
	flagNameValidateSignatureXML        = "signature-xml"
)

// apiValidateSignatureCmd represents the `api validate-signature` command
var apiValidateSignatureCmd = &cobra.Command{
	Use:   "validate-signature",
	Short: "Validate e-factura signature",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fvInvoiceXML, err := cmd.Flags().GetString(flagNameValidateSignatureInvoiceXML)
		if err != nil {
			return err
		}
		fvSignatureXML, err := cmd.Flags().GetString(flagNameValidateSignatureXML)
		if err != nil {
			return err
		}

		client, err := newEfacturaClient(ctx, cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		validateRes, err := client.ValidateSignatureFiles(ctx, fvInvoiceXML, fvSignatureXML)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to validate Invoice UBL XML: %w", err)
		}
		fmt.Println(validateRes.Message)
		return nil
	},
}

func init() {
	apiValidateSignatureCmd.Flags().String(flagNameValidateSignatureInvoiceXML, "", "Path of the Invoice XML file to validate")
	_ = apiValidateSignatureCmd.MarkFlagRequired(flagNameValidateSignatureInvoiceXML)
	apiValidateSignatureCmd.Flags().String(flagNameValidateSignatureXML, "", "Path of the Signature XML file to validate")
	_ = apiValidateSignatureCmd.MarkFlagRequired(flagNameValidateSignatureXML)

	apiCmd.AddCommand(apiValidateSignatureCmd)
}
