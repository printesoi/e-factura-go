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
	"os"

	"github.com/printesoi/e-factura-go/pkg/efactura"
	"github.com/spf13/cobra"
)

const (
	flagNameUploadInvoiceXML        = "invoice-xml"
	flagNameUploadStandard          = "standard"
	flagNameUploadCIF               = "cif"
	flagNameUploadOptionForeign     = "option-foreign"
	flagNameUploadOptionSelfBilled  = "option-self-billed"
	flagNameUploadOptionEnforcement = "option-enforcement"
	flagNameUploadOptionB2C         = "option-b2c"
)

// apiUploadInvoiceCmd represents the `api upload` command
var apiUploadInvoiceCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload an Invoice XML",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fvInvoiceXML, err := cmd.Flags().GetString(flagNameUploadInvoiceXML)
		if err != nil {
			return err
		}

		xmlFile, err := os.Open(fvInvoiceXML)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to read file: %w", err)
		}
		defer xmlFile.Close()

		client, err := newEfacturaClient(ctx, cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvCIF, err := cmd.Flags().GetString(flagNameUploadCIF)
		if err != nil {
			return err
		}

		fvStandard, err := cmd.Flags().GetString(flagNameUploadStandard)
		if err != nil {
			return err
		}
		var standard efactura.UploadStandard
		switch fvStandard {
		case efactura.UploadStandardUBL.String(), "":
			standard = efactura.UploadStandardUBL
		case efactura.UploadStandardCN.String():
			standard = efactura.UploadStandardCN
		case efactura.UploadStandardCII.String():
			standard = efactura.UploadStandardCII
		case efactura.UploadStandardRASP.String():
			standard = efactura.UploadStandardRASP
		default:
			return fmt.Errorf("invalid standard %q", fvStandard)
		}

		var opts []efactura.UploadOption
		if fvB2C, err := cmd.Flags().GetBool(flagNameUploadOptionB2C); err != nil {
			return err
		} else if fvB2C {
			opts = append(opts, efactura.UploadOptionB2C())
		}
		if fvForeign, err := cmd.Flags().GetBool(flagNameUploadOptionForeign); err != nil {
			return err
		} else if fvForeign {
			opts = append(opts, efactura.UploadOptionForeign())
		}
		if fvSelfBilled, err := cmd.Flags().GetBool(flagNameUploadOptionSelfBilled); err != nil {
			return err
		} else if fvSelfBilled {
			opts = append(opts, efactura.UploadOptionSelfBilled())
		}
		if fvEnforcement, err := cmd.Flags().GetBool(flagNameUploadOptionEnforcement); err != nil {
			return err
		} else if fvEnforcement {
			opts = append(opts, efactura.UploadOptionEnforcement())
		}

		resp, err := client.UploadXML(ctx, xmlFile, standard, fvCIF, opts...)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		if !resp.IsOk() {
			cmd.SilenceUsage = true
			return fmt.Errorf("Failed to upload: %s", resp.GetFirstErrorMessage())
		}

		fmt.Printf("Upload OK: upload_index=%d\n", resp.GetUploadIndex())
		return nil
	},
}

func init() {
	apiUploadInvoiceCmd.Flags().String(flagNameUploadInvoiceXML, "", "Path of the Invoice XML file to upload")
	_ = apiUploadInvoiceCmd.MarkFlagRequired(flagNameUploadInvoiceXML)
	apiUploadInvoiceCmd.Flags().String(flagNameUploadCIF, "", "CIF/CUI")
	_ = apiUploadInvoiceCmd.MarkFlagRequired(flagNameUploadCIF)
	apiUploadInvoiceCmd.Flags().String(flagNameUploadStandard, "", "Upload standard (default: UBL)")
	apiUploadInvoiceCmd.Flags().Bool(flagNameUploadOptionB2C, false, "Option specifying that the invoice is B2C")
	apiUploadInvoiceCmd.Flags().Bool(flagNameUploadOptionForeign, false, "Option specifying that the buyer is not a Romanian entity (no CUI or NIF)")
	apiUploadInvoiceCmd.Flags().Bool(flagNameUploadOptionSelfBilled, false, "Option specifying that it's a self-billed invoice")
	apiUploadInvoiceCmd.Flags().Bool(flagNameUploadOptionEnforcement, false, "Option specifying that the invoice is uploaded by the enforcement authority on behalf of the debtor")

	apiCmd.AddCommand(apiUploadInvoiceCmd)
}
