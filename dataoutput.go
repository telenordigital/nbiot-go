package horde

import (
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"
)

type OutputDataMessage struct {
	Device   Device `json:"device"`
	Payload  []byte `json:"payload"`
	Received uint32 `json:"received"`
}

// CollectionOutput calls handler in its own goroutine for each device message received.
// It blocks until an error occurs, including EOF when the connection is closed.
func (c *Client) CollectionOutput(collectionID string, handler func(OutputDataMessage)) error {
	return c.output(fmt.Sprintf("/collections/%s", collectionID), handler)
}

// DeviceOutput calls handler in its own goroutine for each device message received.
// It blocks until an error occurs, including EOF when the connection is closed.
func (c *Client) DeviceOutput(collectionID, deviceID string, handler func(OutputDataMessage)) error {
	return c.output(fmt.Sprintf("/collections/%s/devices/%s", collectionID, deviceID), handler)
}

func (c *Client) output(path string, handler func(OutputDataMessage)) error {
	url, err := url.Parse(c.addr)
	if err != nil {
		return err
	}

	scheme := "wss"
	if url.Scheme == "http" {
		scheme = "ws"
	}

	wscfg, err := websocket.NewConfig(fmt.Sprintf("%s://%s%s/from", scheme, url.Host, path), "http://example.com")
	if err != nil {
		return err
	}
	wscfg.Header.Set("X-API-Token", c.token)

	ws, err := websocket.DialConfig(wscfg)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		var msg struct {
			KeepAlive bool `json:"keepAlive"`
			OutputDataMessage
		}
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			return err
		}

		if msg.KeepAlive {
			continue
		}

		go handler(msg.OutputDataMessage)
	}
}
