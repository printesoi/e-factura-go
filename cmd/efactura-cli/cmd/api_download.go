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
	"io"
	"os"

	"github.com/spf13/cobra"
)

const (
	flagNameDownloadID      = "id"
	flagNameDownloadOutFile = "out"
)

// apiDownloadCmd represents the `api download` command
var apiDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download an Invoice zip",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := newEfacturaClient(ctx, cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvDownloadID, err := cmd.Flags().GetInt64(flagNameDownloadID)
		if err != nil {
			return err
		}

		res, err := client.DownloadInvoice(ctx, fvDownloadID)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		if !res.IsOk() {
			cmd.SilenceUsage = true
			return fmt.Errorf("error downloading invoice %d: %s", fvDownloadID, res.Error.Error)
		}

		fvOutFile, err := cmd.Flags().GetString(flagNameDownloadOutFile)
		if err != nil {
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
		if _, err := io.Copy(out, bytes.NewReader(res.Zip)); err != nil {
			cmd.SilenceUsage = true
			fmt.Fprintf(os.Stderr, "error copying data for invoice %d: %v\n", fvDownloadID, err)
			return err
		}

		return nil
	},
}

func init() {
	apiDownloadCmd.Flags().Int64(flagNameDownloadID, 0, "Download ID for the Invoice zip")
	_ = apiDownloadCmd.MarkFlagRequired(flagNameDownloadID)

	apiDownloadCmd.Flags().String(flagNameDownloadOutFile, "", "Write output to this file instead of stdout")

	apiCmd.AddCommand(apiDownloadCmd)
}
