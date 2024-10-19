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

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"context"
	"errors"
	"net/url"
	"sync"

	xoauth2 "golang.org/x/oauth2"

	ioauth2 "github.com/printesoi/e-factura-go/internal/oauth2"
)

// TokenChangedHandler is a handler provided to
// Config.TokenSourceWithChangedHandler or Config.TokenRefresher that is called
// when the token changes (is refreshed).
type TokenChangedHandler func(ctx context.Context, t *xoauth2.Token) error

// tokenFromInternal maps an *ioauth2.Token struct into
// a *Token struct.
func tokenFromInternal(t *ioauth2.Token) *xoauth2.Token {
	if t == nil {
		return nil
	}
	return &xoauth2.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
	}
}

// retrieveToken takes a *Config and uses that to retrieve an *ioauth2.Token.
// This token is then mapped from *ioauth2.Token into an *oauth2.Token which
// is returned along with an error..
func retrieveToken(ctx context.Context, c *xoauth2.Config, v url.Values) (*xoauth2.Token, error) {
	tk, err := ioauth2.RetrieveToken(ctx, c.ClientID, c.ClientSecret, c.Endpoint.TokenURL, v, ioauth2.AuthStyle(c.Endpoint.AuthStyle), nil)
	if err != nil {
		if rErr, ok := err.(*ioauth2.RetrieveError); ok {
			return nil, (*xoauth2.RetrieveError)(rErr)
		}
		return nil, err
	}
	return tokenFromInternal(tk), nil
}

// tokenRefresher is a TokenSource that makes "grant_type"=="refresh_token"
// HTTP requests to renew a token using a RefreshToken. When the token is
// refreshed, onTokenChanged is called.
type tokenRefresher struct {
	ctx            context.Context // used to get HTTP requests
	conf           *xoauth2.Config
	refreshToken   string
	onTokenChanged TokenChangedHandler
}

// WARNING: Token is not safe for concurrent access, as it
// updates the tokenRefresher's refreshToken field.
// Within this package, it is used by reuseTokenSource which
// synchronizes calls to this method with its own mutex.
func (tf *tokenRefresher) Token() (*xoauth2.Token, error) {
	if tf.refreshToken == "" {
		return nil, errors.New("oauth2: token expired and refresh token is not set")
	}

	tk, err := retrieveToken(tf.ctx, tf.conf, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {tf.refreshToken},
	})

	if err != nil {
		return nil, err
	}
	if tf.refreshToken != tk.RefreshToken {
		tf.refreshToken = tk.RefreshToken
		if tf.onTokenChanged != nil {
			if err := tf.onTokenChanged(tf.ctx, tk); err != nil {
				return tk, err
			}
		}
	}
	return tk, err
}

// reuseTokenSource is a TokenSource that holds a single token in memory
// and validates its expiry before each call to retrieve it with
// Token. If it's expired, it will be auto-refreshed using the
// new TokenSource.
type reuseTokenSource struct {
	new xoauth2.TokenSource // called when t is expired.

	mu sync.Mutex // guards t
	t  *xoauth2.Token
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
func (s *reuseTokenSource) Token() (*xoauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}

// TokenRefresher returns a TokenSource that makes "grant_type"=="refresh_token"
// HTTP requests to renew a token using a RefreshToken.
// WARNING: the returned TokenSource is not safe for concurrent access, so you
// need to protect it with a mutex. It's recommended to use TokenSource instead.
func (c *Config) TokenRefresher(ctx context.Context, t *xoauth2.Token, onTokenChanged TokenChangedHandler) xoauth2.TokenSource {
	if t == nil || t.RefreshToken == "" {
		return nil
	}
	return &tokenRefresher{
		ctx:            ctx,
		conf:           &c.Config,
		refreshToken:   t.RefreshToken,
		onTokenChanged: onTokenChanged,
	}
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided context. The
// returned TokenSource is safe for concurrent access.
func (c *Config) TokenSource(ctx context.Context, t *xoauth2.Token) xoauth2.TokenSource {
	return c.TokenSourceWithChangedHandler(ctx, t, nil)
}

// TokenSourceWithChangedHandler returns a TokenSource that returns t until t
// expires, automatically refreshing it as necessary using the provided
// context. Every time the access token is refreshed, the onTokenChanged
// handler is called. This is useful if you need to update the token in a
// store/db. The returned TokenSource is safe for concurrent access.
func (c *Config) TokenSourceWithChangedHandler(ctx context.Context, t *xoauth2.Token, onTokenChanged TokenChangedHandler) xoauth2.TokenSource {
	tkr := &tokenRefresher{
		ctx:            ctx,
		conf:           &c.Config,
		onTokenChanged: onTokenChanged,
	}
	if t != nil {
		tkr.refreshToken = t.RefreshToken
	}
	return &reuseTokenSource{
		t:   t,
		new: tkr,
	}
}
