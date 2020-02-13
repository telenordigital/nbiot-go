package nbiot

import (
	"fmt"
	"net/http"
)

// DownstreamMessage is a message to be sent to a device.
type DownstreamMessage struct {
	Port      int    `json:"port"`
	Payload   []byte `json:"payload"`
	Path      string `json:"coapPath,omitempty"`  // This is used by CoAP support to specify the path
	Transport string `json:"transport,omitempty"` // This is used by CoAP support to specify transport
}

// Send sends a message to a device.
func (c *Client) Send(collectionID, deviceID string, msg DownstreamMessage) error {
	return c.request(http.MethodPost, fmt.Sprintf("/collections/%s/devices/%s/to", collectionID, deviceID), msg, nil)
}

// Broadcast sends a message to all devices in a collection.
func (c *Client) Broadcast(collectionID string, msg DownstreamMessage) (result BroadcastResult, err error) {
	err = c.request(http.MethodPost, fmt.Sprintf("/collections/%s/to", collectionID), msg, &result)
	return result, err
}

// BroadcastResult is the result of a broadcast.
type BroadcastResult struct {
	Sent   int              `json:"sent"`
	Failed int              `json:"failed"`
	Errors []BroadcastError `json:"errors"`
}

// BroadcastError is an error from a broadcast.
type BroadcastError struct {
	DeviceID string `json:"deviceId"`
	Message  string `json:"message"`
}
