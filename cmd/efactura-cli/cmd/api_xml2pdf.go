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

	"github.com/printesoi/e-factura-go/pkg/efactura"
	"github.com/spf13/cobra"
)

const (
	flagNameXML2PDFInvoiceXML = "xml"
	flagNameXML2PDFOutFile    = "out"
	flagNameXML2PDFNoValidate = "no-validate"
)

// apiXML2PDFCmd represents the `api download` command
var apiXML2PDFCmd = &cobra.Command{
	Use:   "xml2pdf",
	Short: "Convert XML to PDF",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := newEfacturaClient(ctx, cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvInvoiceXML, err := cmd.Flags().GetString(flagNameXML2PDFInvoiceXML)
		if err != nil {
			return err
		}

		xmlFile, err := os.Open(fvInvoiceXML)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		defer xmlFile.Close()

		fvNoValidate, err := cmd.Flags().GetBool(flagNameXML2PDFNoValidate)
		if err != nil {
			return err
		}

		fvOutFile, err := cmd.Flags().GetString(flagNameXML2PDFOutFile)
		if err != nil {
			return err
		}

		res, err := client.XMLToPDF(ctx, xmlFile, efactura.ValidateStandardFACT1, fvNoValidate)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		out := os.Stdout
		if fvOutFile != "" {
			out, err = os.OpenFile(fvOutFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			defer func() {
				if err := out.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "error closing file: %v\n", err)
				}
			}()
		}
		if _, err := out.Write(res.PDF); err != nil {
			cmd.SilenceUsage = true
			fmt.Fprintf(os.Stderr, "error writing PDF data: %v\n", err)
			return err
		}

		return nil
	},
}

func init() {
	apiXML2PDFCmd.Flags().String(flagNameXML2PDFInvoiceXML, "", "Input Invoice XML file path")
	_ = apiXML2PDFCmd.MarkFlagRequired(flagNameXML2PDFInvoiceXML)

	apiXML2PDFCmd.Flags().String(flagNameXML2PDFOutFile, "", "Write output to this file instead of stdout")

	apiXML2PDFCmd.Flags().Bool(flagNameXML2PDFNoValidate, false, "Skip validation of XML file")

	apiCmd.AddCommand(apiXML2PDFCmd)
}
