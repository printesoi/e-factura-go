// Copyright 2026 Victor Dodon
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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

const (
	flagNameCallbackPath = "callback-path"
	flagNameBindAddr     = "addr"
)

// authAuthorizeServerCmd represents the `auth authorize-server` command
var authAuthorizeServerCmd = &cobra.Command{
	Use:   "authorize-server",
	Short: "Exchange an OAuth device code for a token",
	Long:  `Exchange an OAuth device code for a token`,
	RunE: func(cmd *cobra.Command, args []string) error {
		oauth2Cfg, err := NewOAuth2Config(cmd)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		bindAddr, err := cmd.Flags().GetString(flagNameBindAddr)
		if err != nil {
			return err
		}

		callbackPath, err := cmd.Flags().GetString(flagNameCallbackPath)
		if err != nil {
			return err
		}

		mux := http.NewServeMux()
		mux.HandleFunc(callbackPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.RequestURI)

			code := r.URL.Query().Get("code")
			if code == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			token, err := oauth2Cfg.Exchange(context.Background(), code)
			if err != nil {
				log.Printf("Failed to exchange code: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			tokenJSON, err := json.Marshal(token)
			if err != nil {
				log.Printf("Failed marshal access token: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Printf("Access token: %s", string(tokenJSON))

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Access token successfully generated, check command output for the access token"))
		}))
		mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.RequestURI)
			w.WriteHeader(http.StatusNotFound)
		}))

		server := &http.Server{
			Addr:    bindAddr,
			Handler: mux,
		}

		// Channel to listen for OS signals
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		go func() {
			log.Printf("Server listening on %s", bindAddr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe error: %v", err)
			}
		}()

		// Wait for shutdown signal
		<-stop
		log.Println("Shutdown signal received")

		// Context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		return nil
	},
}

func init() {
	authAuthorizeServerCmd.Flags().String(flagNameCallbackPath, "/callback", "The callback path")
	authAuthorizeServerCmd.Flags().String(flagNameBindAddr, "localhost:5000", "Bind addr")

	AuthCmd.AddCommand(authAuthorizeServerCmd)
}
