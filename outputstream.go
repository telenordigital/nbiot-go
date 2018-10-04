package nbiot

import (
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"
)

type OutputStream struct {
	ws *websocket.Conn
}

type OutputDataMessage struct {
	Device   Device `json:"device"`
	Payload  []byte `json:"payload"`
	Received int64  `json:"received"`
}

func (c *Client) CollectionOutputStream(collectionID string) (*OutputStream, error) {
	return c.outputStream(fmt.Sprintf("/collections/%s", collectionID))
}

func (c *Client) DeviceOutputStream(collectionID, deviceID string) (*OutputStream, error) {
	return c.outputStream(fmt.Sprintf("/collections/%s/devices/%s", collectionID, deviceID))
}

func (c *Client) outputStream(path string) (*OutputStream, error) {
	url, err := url.Parse(c.addr)
	if err != nil {
		return nil, err
	}

	scheme := "wss"
	if url.Scheme == "http" {
		scheme = "ws"
	}

	wscfg, err := websocket.NewConfig(fmt.Sprintf("%s://%s%s/from", scheme, url.Host, path), "http://example.com")
	if err != nil {
		return nil, err
	}
	wscfg.Header.Set("X-API-Token", c.token)

	ws, err := websocket.DialConfig(wscfg)
	if err != nil {
		return nil, err
	}

	return &OutputStream{ws}, nil
}

func (s *OutputStream) Recv() (OutputDataMessage, error) {
	for {
		var msg struct {
			Type string `json:"type"`
			OutputDataMessage
		}
		err := websocket.JSON.Receive(s.ws, &msg)
		if err != nil {
			return OutputDataMessage{}, err
		}

		if msg.Type == "data" {
			return msg.OutputDataMessage, nil
		}
	}
}

func (s *OutputStream) Close() {
	s.ws.Close()
}
