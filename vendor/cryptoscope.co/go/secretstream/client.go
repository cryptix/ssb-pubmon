/*
This file is part of secretstream.

secretstream is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

secretstream is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with secretstream.  If not, see <http://www.gnu.org/licenses/>.
*/

package secretstream // import "cryptoscope.co/go/secretstream"

import (
	"context"
	"net"
	"strings"
	"time"

	"cryptoscope.co/go/secretstream/boxstream"
	"cryptoscope.co/go/secretstream/secrethandshake"
	"github.com/agl/ed25519"
)

// Client can dial secret-handshake server endpoints
type Client struct {
	appKey []byte
	kp     secrethandshake.EdKeyPair
}

// NewClient creates a new Client with the passed keyPair and appKey
func NewClient(kp secrethandshake.EdKeyPair, appKey []byte) (*Client, error) {
	// TODO: consistancy check?!..
	return &Client{
		appKey: appKey,
		kp:     kp,
	}, nil
}

// NewDialer returns a net.Dial-esque dialer that does a secrethandshake key exchange
// and wraps the underlying connection into a boxstream
func (c *Client) NewDialer(pubKey [ed25519.PublicKeySize]byte) (Dialer, error) {
	var d net.Dialer
	return func(ctx context.Context, n, a string) (net.Conn, error) {
		if !strings.HasPrefix(n, "tcp") {
			return nil, ErrOnlyTCP
		}
		conn, err := d.DialContext(ctx, n, a)
		if err != nil {
			return nil, err
		}
		// TODO: tune
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		state, err := secrethandshake.NewClientState(c.appKey, c.kp, pubKey)
		if err != nil {
			return nil, err
		}

		if err := secrethandshake.Client(state, conn); err != nil {
			return nil, err
		}

		enKey, enNonce := state.GetBoxstreamEncKeys()
		deKey, deNonce := state.GetBoxstreamDecKeys()

		conn.SetDeadline(time.Now().Add(time.Hour))
		boxed := Conn{
			Reader: boxstream.NewUnboxer(conn, &deNonce, &deKey),
			Writer: boxstream.NewBoxer(conn, &enNonce, &enKey),
			conn:   conn,
			local:  c.kp.Public[:],
			remote: state.Remote(),
		}

		return boxed, nil
	}, nil
}
