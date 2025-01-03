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

	"github.com/spf13/cobra"
)

const (
	flagNameValidateSignatureZip = "zip"
)

// apiValidateSignatureZipCmd represents the `api download` command
var apiValidateSignatureZipCmd = &cobra.Command{
	Use:   "validate-signature-zip",
	Short: "Validate e-factura signature ZIP archive",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fvInvoiceZip, err := cmd.Flags().GetString(flagNameValidateSignatureZip)
		if err != nil {
			return err
		}

		client, err := newEfacturaClient(ctx, cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		validateRes, err := client.ValidateSignatureZipFile(ctx, fvInvoiceZip)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to validate Invoice ZIP: %w", err)
		}
		fmt.Println(validateRes.Message)
		return nil
	},
}

func init() {
	apiValidateSignatureZipCmd.Flags().String(flagNameValidateSignatureZip, "", "Path to the ZIP archive")
	_ = apiValidateSignatureZipCmd.MarkFlagRequired(flagNameValidateSignatureZip)

	apiCmd.AddCommand(apiValidateSignatureZipCmd)
}
