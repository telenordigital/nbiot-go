package nbiot

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// OutputStream provides a stream of OutputDataMessages.
type OutputStream struct {
	ws *websocket.Conn
}

// OutputDataMessage represents a message sent by a device.
type OutputDataMessage struct {
	Device       Device `json:"device"`
	Payload      []byte `json:"payload"`
	Received     int64  `json:"received"`
	Transport    string `json:"transport"`
	CoAPMetaData struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	} `json:"coapMetaData"`
	UDPMetaData struct {
		LocalPort  int `json:"localPort"`
		RemotePort int `json:"remotePort"`
	} `json:"udpMetaData"`
}

// CollectionOutputStream streams messages from all devices in a collection.
func (c *Client) CollectionOutputStream(collectionID string) (*OutputStream, error) {
	return c.outputStream(fmt.Sprintf("/collections/%s", collectionID))
}

// DeviceOutputStream streams messages from one device.
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

	urlStr := fmt.Sprintf("%s://%s%s/from", scheme, url.Host, path)

	header := http.Header{}
	header.Add("X-API-Token", c.token)

	dialer := websocket.Dialer{}
	ws, _, err := dialer.Dial(urlStr, header)
	if err != nil {
		return nil, err
	}

	return &OutputStream{ws}, nil
}

// Recv blocks until a new message is received.
// It returns io.EOF if the stream is closed by the server.
func (s *OutputStream) Recv() (OutputDataMessage, error) {
	for {
		var msg struct {
			Type string `json:"type"`
			OutputDataMessage
		}
		err := s.ws.ReadJSON(&msg)
		if err != nil {
			return OutputDataMessage{}, err
		}

		if msg.Type == "data" {
			return msg.OutputDataMessage, nil
		}
	}
}

// Close closes the output stream.
func (s *OutputStream) Close() {
	s.ws.Close()
}
