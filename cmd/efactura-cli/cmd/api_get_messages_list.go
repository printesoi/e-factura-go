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
	"encoding/json"
	"fmt"
	"time"

	"github.com/printesoi/e-factura-go/efactura"
	"github.com/spf13/cobra"
)

const (
	flagNameMessageFilter  = "filter"
	flagNameMessageNumDays = "num-days"
	flagNameMessageCIF     = "cif"
	flagNameMessageJSON    = "json"
)

// apiGetMessagesCmd represents the `api get-messages` command
var apiGetMessagesCmd = &cobra.Command{
	Use:   "get-messages",
	Short: "Get messages list",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := newEfacturaClient(ctx, cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		fvCIF, err := cmd.Flags().GetString(flagNameMessageCIF)
		if err != nil {
			return err
		}

		fvNumDays, err := cmd.Flags().GetInt(flagNameMessageNumDays)
		if err != nil {
			return err
		}

		fvFilter, err := cmd.Flags().GetString(flagNameMessageFilter)
		if err != nil {
			return err
		}

		messageFilter := efactura.MessageFilterAll
		switch fvFilter {
		case efactura.MessageFilterErrors.String():
			messageFilter = efactura.MessageFilterErrors
		case efactura.MessageFilterSent.String():
			messageFilter = efactura.MessageFilterSent
		case efactura.MessageFilterReceived.String():
			messageFilter = efactura.MessageFilterReceived
		case efactura.MessageFilterBuyerMessage.String():
			messageFilter = efactura.MessageFilterBuyerMessage
		case "A", "":
			messageFilter = efactura.MessageFilterAll
		default:
			return fmt.Errorf("invalid filter '%s'", fvFilter)
		}
		res, err := client.GetMessagesList(ctx, fvCIF, fvNumDays, messageFilter)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		if !res.IsOk() {
			cmd.SilenceUsage = true
			return fmt.Errorf("error getting messages list: %s", res.Error)
		}

		asJson, err := cmd.Flags().GetBool(flagNameMessageJSON)
		if err != nil {
			return err
		}
		if asJson {
			messagesData, err := json.Marshal(res)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			fmt.Printf("%s\n", string(messagesData))
		} else {
			fmt.Printf("Got %d messages:\n", len(res.Messages))
			for i, m := range res.Messages {
				t, _ := time.Parse("200601021504", m.CreationDate)
				fmt.Printf("%d. [%s] CIF=%s, Type=%s, Upload_index=%d, ID=%d, Details=%s\n",
					i, t.Format(time.DateTime), m.CIF, m.Type, m.GetUploadIndex(), m.GetID(), m.Details)
			}
		}

		return nil
	},
}

func init() {
	apiGetMessagesCmd.Flags().Int(flagNameMessageNumDays, 1, "The number of days to fetch messages (must be >=1, <= 60)")
	_ = apiGetMessagesCmd.MarkFlagRequired(flagNameDownloadID)

	apiGetMessagesCmd.Flags().String(flagNameMessageFilter, "A", "Message filter")

	apiGetMessagesCmd.Flags().String(flagNameMessageCIF, "", "CIF/CUI")
	_ = apiGetMessagesCmd.MarkFlagRequired(flagNameMessageCIF)

	apiGetMessagesCmd.Flags().Bool(flagNameMessageJSON, false, "Output RAW json")

	apiCmd.AddCommand(apiGetMessagesCmd)
}
