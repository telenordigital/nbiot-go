package nbiot

import (
	"fmt"
	"time"
)

// Device represents a device.
type Device struct {
	DeviceID     *string           `json:"deviceId"`
	CollectionID *string           `json:"collectionId,omitempty"`
	IMEI         *string           `json:"imei"`
	IMSI         *string           `json:"imsi"`
	Tags         map[string]string `json:"tags,omitempty"`
}

// Device gets a device.
func (c *Client) Device(collectionID, deviceID string) (Device, error) {
	var device Device
	err := c.get(fmt.Sprintf("/collections/%s/devices/%s", collectionID, deviceID), &device)
	return device, err
}

// Devices gets all devices in the collection.
func (c *Client) Devices(collectionID string) ([]Device, error) {
	var devices struct {
		Devices []Device `json:"devices"`
	}
	err := c.get(fmt.Sprintf("/collections/%s/devices", collectionID), &devices)
	return devices.Devices, err
}

// CreateDevice creates a device.
func (c *Client) CreateDevice(collectionID string, device Device) (Device, error) {
	err := c.create(fmt.Sprintf("/collections/%s/devices", collectionID), &device)
	return device, err
}

// UpdateDevice updates a device.
// No tags are deleted, only added or updated.
func (c *Client) UpdateDevice(collectionID string, device Device) (Device, error) {
	err := c.update(fmt.Sprintf("/collections/%s/devices/%s", collectionID, *device.DeviceID), &device)
	return device, err
}

// DeleteDeviceTag deletes a tag from a device.
func (c *Client) DeleteDeviceTag(collectionID, deviceID, name string) error {
	return c.delete(fmt.Sprintf("/collections/%s/devices/%s/tags/%s", collectionID, deviceID, name))
}

// DeleteDevice deletes a device.
func (c *Client) DeleteDevice(collectionID, deviceID string) error {
	return c.delete(fmt.Sprintf("/collections/%s/devices/%s", collectionID, deviceID))
}

// DeviceData returns all the stored data for the device.
func (c *Client) DeviceData(collectionID, deviceID string, since time.Time, until time.Time, limit int) ([]Datapoint, error) {
	var s, u int64
	if !since.IsZero() {
		s = since.UnixNano() / int64(time.Millisecond)
	}
	if !until.IsZero() {
		u = until.UnixNano() / int64(time.Millisecond)
	}

	var data struct {
		Datapoints []Datapoint `json:"messages"`
	}
	err := c.get(fmt.Sprintf("/collections/%s/devices/%s/data?since=%d&until=%d&limit=%d", collectionID, deviceID, s, u, limit), &data)
	return data.Datapoints, err
}
