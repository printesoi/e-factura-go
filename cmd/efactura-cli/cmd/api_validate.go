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
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/printesoi/e-factura-go/pkg/efactura"
	"github.com/spf13/cobra"
)

const (
	flagNameValidateInvoiceXML = "invoice-xml"
)

// apiValidateInvoiceCmd represents the `api validate-invoice` command
var apiValidateInvoiceCmd = &cobra.Command{
	Use:   "validate-invoice",
	Short: "Validate Invoice UBL XML",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fvInvoiceXML, err := cmd.Flags().GetString(flagNameValidateInvoiceXML)
		if err != nil {
			return err
		}

		xmlData, err := os.ReadFile(fvInvoiceXML)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to read file: %w", err)
		}

		client, err := newEfacturaClient(ctx, cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		validateRes, err := client.ValidateXML(ctx, bytes.NewReader(xmlData), efactura.ValidateStandardFACT1)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to validate Invoice UBL XML: %w", err)
		}
		if validateRes.IsOk() {
			fmt.Println("successfully validated Invoice UBL XML")
		} else {
			cmd.SilenceUsage = true
			return fmt.Errorf("validate Invoice UBL XML failed: %s", validateRes.GetFirstMessage())
		}

		return nil
	},
}

func init() {
	apiValidateInvoiceCmd.Flags().String(flagNameValidateInvoiceXML, "", "Path of the Invoice XML file to validate")
	_ = apiValidateInvoiceCmd.MarkFlagRequired(flagNameValidateInvoiceXML)

	apiCmd.AddCommand(apiValidateInvoiceCmd)
}
