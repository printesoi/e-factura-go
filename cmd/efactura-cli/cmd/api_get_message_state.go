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
	flagNameGetMessageStateUploadIndex = "upload-index"
)

// apiGetMessageStateCmd represents the `api download` command
var apiGetMessageStateCmd = &cobra.Command{
	Use:   "get-message-state",
	Short: "Get message state for an upload index",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := newEfacturaClient(ctx, cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvUploadIndex, err := cmd.Flags().GetInt64(flagNameGetMessageStateUploadIndex)
		if err != nil {
			return err
		}

		res, err := client.GetMessageState(ctx, fvUploadIndex)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("get message state failed: %w", err)
		}
		if res.IsOk() {
			fmt.Printf("Message for upload index %d: ok, download_id=%d\n", fvUploadIndex, res.GetDownloadID())
		} else if res.IsProcessing() {
			fmt.Printf("Message for upload index %d: processing\n", fvUploadIndex)
		} else if res.IsInvalidXML() {
			fmt.Printf("Message for upload index %d: invalid XML, download_id=%d, message=%s\n", fvUploadIndex, res.GetDownloadID(), res.GetFirstErrorMessage())
		} else if res.IsNok() {
			fmt.Printf("Message for upload index %d: nok, download_id=%d, message=%s\n", fvUploadIndex, res.GetDownloadID(), res.GetFirstErrorMessage())
		} else {
			fmt.Printf("Message for upload index %d: unknown state '%s', message: '%s'\n", fvUploadIndex, res.State, res.GetFirstErrorMessage())
		}

		return nil
	},
}

func init() {
	apiGetMessageStateCmd.Flags().Int64(flagNameGetMessageStateUploadIndex, 0, "Upload index")
	_ = apiGetMessageStateCmd.MarkFlagRequired(flagNameGetMessageStateUploadIndex)

	apiCmd.AddCommand(apiGetMessageStateCmd)
}
