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
	flagNameXml2PdfInFilePath  = "xml-file"
	flagNameXml2PdfOutFilePath = "out"
	flagNameXml2PdfStandard    = "standard"
	flagNameXml2PdfNoValidate  = "no-validate"
)

// publicApiXml2PdfCmd represents the `public-api xml2pdf` command
var publicApiXml2PdfCmd = &cobra.Command{
	Use:   "xml2pdf",
	Short: "Convert XML to PDF",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newEfacturaPublicClient(cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvXmlFile, err := cmd.Flags().GetString(flagNameXml2PdfInFilePath)
		if err != nil {
			return err
		}
		f, err := os.Open(fvXmlFile)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		defer f.Close()

		fvOutFile, err := cmd.Flags().GetString(flagNameXml2PdfOutFilePath)
		if err != nil {
			return err
		}

		fvStandard, err := cmd.Flags().GetString(flagNameValidateStandard)
		if err != nil {
			return err
		}
		switch fvStandard {
		case string(efactura.ValidateStandardFACT1), string(efactura.ValidateStandardFCN):
		default:
			return fmt.Errorf("invalid standard: `%s`", fvStandard)
		}

		fvNoValidate, err := cmd.Flags().GetBool(flagNameXml2PdfNoValidate)
		if err != nil {
			return err
		}

		xml2PdfRes, err := client.XMLToPDF(context.Background(), f, efactura.ValidateStandard(fvStandard), fvNoValidate)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		if !xml2PdfRes.IsOk() {
			cmd.SilenceUsage = true
			return fmt.Errorf("XML-to-PDF error: %s", xml2PdfRes.GetError().GetFirstMessage())
		}

		if err := os.WriteFile(fvOutFile, xml2PdfRes.PDF, 0644); err != nil {
			cmd.SilenceUsage = true
			return err
		}
		return nil
	},
}

func init() {
	publicApiXml2PdfCmd.Flags().StringP(flagNameXml2PdfInFilePath, "f", "", "Path of the input XML file")
	_ = publicApiXml2PdfCmd.MarkFlagRequired(flagNameXml2PdfInFilePath)

	publicApiXml2PdfCmd.Flags().StringP(flagNameXml2PdfStandard, "s", "", fmt.Sprintf("Standard (%s/%s)",
		efactura.ValidateStandardFACT1, efactura.ValidateStandardFCN))
	_ = publicApiXml2PdfCmd.MarkFlagRequired(flagNameXml2PdfStandard)

	publicApiXml2PdfCmd.Flags().StringP(flagNameXml2PdfOutFilePath, "o", "", "Path of the output PDF file")
	_ = publicApiXml2PdfCmd.MarkFlagRequired(flagNameXml2PdfOutFilePath)

	publicApiXml2PdfCmd.Flags().Bool(flagNameXml2PdfNoValidate, false, "Don't validate XML")

	publicApiCmd.AddCommand(publicApiXml2PdfCmd)
}
