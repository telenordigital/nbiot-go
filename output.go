package nbiot

import "fmt"

// Output represents a data output for a collection.
// WebHookOutput, MQTTOutput, IFTTTOutput, and UDPOutput implement this interface.
type Output interface {
	GetID() string
	GetCollectionID() string
	IsDisabled() bool
	GetTags() map[string]string

	toOutput() output
}

// WebHookOutput describes a webhook output.
type WebHookOutput struct {
	ID                string
	CollectionID      string
	URL               string
	BasicAuthUser     string
	BasicAuthPass     string
	CustomHeaderName  string
	CustomHeaderValue string
	Disabled          bool
	Tags              map[string]string
}

// MQTTOutput describes an MQTT output.
type MQTTOutput struct {
	ID               string
	CollectionID     string
	Endpoint         string
	DisableCertCheck bool
	Username         string
	Password         string
	ClientID         string
	TopicName        string
	Disabled         bool
	Tags             map[string]string
}

// IFTTTOutput describes an IFTTT output.
type IFTTTOutput struct {
	ID           string
	CollectionID string
	Key          string
	EventName    string
	AsIsPayload  bool
	Disabled     bool
	Tags         map[string]string
}

// UDPOutput describes a UDP output.
type UDPOutput struct {
	ID           string
	CollectionID string
	Host         string
	Port         int
	Disabled     bool
	Tags         map[string]string
}

// Output retrieves an output
func (c *Client) Output(collectionID, outputID string) (Output, error) {
	var output output
	err := c.get(fmt.Sprintf("/collections/%s/outputs/%s", collectionID, outputID), &output)
	if err != nil {
		return nil, err
	}
	return output.toOutput()
}

// Outputs retrieves a list of outputs on a collection
func (c *Client) Outputs(collectionID string) ([]Output, error) {
	var outputs struct {
		Outputs []output `json:"outputs"`
	}
	err := c.get(fmt.Sprintf("/collections/%s/outputs", collectionID), &outputs)
	if err != nil {
		return nil, err
	}

	ret := make([]Output, len(outputs.Outputs))
	for i, o := range outputs.Outputs {
		ret[i], err = o.toOutput()
		if err != nil {
			return nil, err
		}
	}
	return ret, err
}

// CreateOutput creates an output
func (c *Client) CreateOutput(collectionID string, output Output) (Output, error) {
	o := output.toOutput()
	err := c.create(fmt.Sprintf("/collections/%s/outputs", collectionID), &o)
	if err != nil {
		return nil, err
	}
	return o.toOutput()
}

// UpdateOutput updates an output. The type field can't be modified
// No tags are deleted, only added or updated.
func (c *Client) UpdateOutput(collectionID string, output Output) (Output, error) {
	o := output.toOutput()
	err := c.update(fmt.Sprintf("/collections/%s/outputs/%s", collectionID, *o.ID), &o)
	if err != nil {
		return nil, err
	}
	return o.toOutput()
}

// OutputLogEntry is an entry in an output log.
type OutputLogEntry struct {
	Message   string `json:"message"`   // The message itself
	Timestamp int64  `json:"timestamp"` // Time in ms
	Repeated  uint8  `json:"repeated"`  // Repeat count
}

// OutputLogs returns the logs for an output.
func (c *Client) OutputLogs(collectionID, outputID string) ([]OutputLogEntry, error) {
	var log struct {
		Logs []OutputLogEntry `json:"logs"`
	}
	err := c.get(fmt.Sprintf("/collections/%s/outputs/%s/logs", collectionID, outputID), &log)
	return log.Logs, err
}

// OutputStatus is the status of an output.
type OutputStatus struct {
	ErrorCount int `json:"errorCount"`
	Forwarded  int `json:"forwarded"`
	Received   int `json:"received"`
	Retries    int `json:"retries"`
}

// OutputStatus returns the status for an output.
func (c *Client) OutputStatus(collectionID, outputID string) (stat OutputStatus, err error) {
	err = c.get(fmt.Sprintf("/collections/%s/outputs/%s/status", collectionID, outputID), &stat)
	return stat, err
}

// DeleteOutputTag deletes a tag from an output.
func (c *Client) DeleteOutputTag(collectionID, outputID, name string) error {
	return c.delete(fmt.Sprintf("/collections/%s/outputs/%s/tags/%s", collectionID, outputID, name))
}

// DeleteOutput removes an output
func (c *Client) DeleteOutput(collectionID, outputID string) error {
	return c.delete(fmt.Sprintf("/collections/%s/outputs/%s", collectionID, outputID))
}

// GetID returns the output ID.
func (o WebHookOutput) GetID() string { return o.ID }

// GetID returns the output ID.
func (o MQTTOutput) GetID() string { return o.ID }

// GetID returns the output ID.
func (o IFTTTOutput) GetID() string { return o.ID }

// GetID returns the output ID.
func (o UDPOutput) GetID() string { return o.ID }

// GetCollectionID returns the collection ID.
func (o WebHookOutput) GetCollectionID() string { return o.CollectionID }

// GetCollectionID returns the collection ID.
func (o MQTTOutput) GetCollectionID() string { return o.CollectionID }

// GetCollectionID returns the collection ID.
func (o IFTTTOutput) GetCollectionID() string { return o.CollectionID }

// GetCollectionID returns the collection ID.
func (o UDPOutput) GetCollectionID() string { return o.CollectionID }

// IsDisabled returns whether the output is disabled.
func (o WebHookOutput) IsDisabled() bool { return o.Disabled }

// IsDisabled returns whether the output is disabled.
func (o MQTTOutput) IsDisabled() bool { return o.Disabled }

// IsDisabled returns whether the output is disabled.
func (o IFTTTOutput) IsDisabled() bool { return o.Disabled }

// IsDisabled returns whether the output is disabled.
func (o UDPOutput) IsDisabled() bool { return o.Disabled }

// GetTags returns the output's tags.
func (o WebHookOutput) GetTags() map[string]string { return o.Tags }

// GetTags returns the output's tags.
func (o MQTTOutput) GetTags() map[string]string { return o.Tags }

// GetTags returns the output's tags.
func (o IFTTTOutput) GetTags() map[string]string { return o.Tags }

// GetTags returns the output's tags.
func (o UDPOutput) GetTags() map[string]string { return o.Tags }

func (o WebHookOutput) toOutput() output {
	typ := "webhook"
	enabled := !o.Disabled
	return output{
		ID:           &o.ID,
		CollectionID: &o.CollectionID,
		Type:         &typ,
		Config: map[string]interface{}{
			"url":               o.URL,
			"basicAuthUser":     o.BasicAuthUser,
			"basicAuthPass":     o.BasicAuthPass,
			"customHeaderName":  o.CustomHeaderName,
			"customHeaderValue": o.CustomHeaderValue,
		},
		Enabled: &enabled,
		Tags:    o.Tags,
	}
}

func (o MQTTOutput) toOutput() output {
	typ := "mqtt"
	enabled := !o.Disabled
	return output{
		ID:           &o.ID,
		CollectionID: &o.CollectionID,
		Type:         &typ,
		Config: map[string]interface{}{
			"endpoint":         o.Endpoint,
			"disableCertCheck": o.DisableCertCheck,
			"username":         o.Username,
			"password":         o.Password,
			"clientId":         o.ClientID,
			"topicName":        o.TopicName,
		},
		Enabled: &enabled,
		Tags:    o.Tags,
	}
}

func (o IFTTTOutput) toOutput() output {
	typ := "ifttt"
	enabled := !o.Disabled
	return output{
		ID:           &o.ID,
		CollectionID: &o.CollectionID,
		Type:         &typ,
		Config: map[string]interface{}{
			"key":         o.Key,
			"eventName":   o.EventName,
			"asIsPayload": o.AsIsPayload,
		},
		Enabled: &enabled,
		Tags:    o.Tags,
	}
}

func (o UDPOutput) toOutput() output {
	typ := "udp"
	enabled := !o.Disabled
	return output{
		ID:           &o.ID,
		CollectionID: &o.CollectionID,
		Type:         &typ,
		Config: map[string]interface{}{
			"host": o.Host,
			"port": o.Port,
		},
		Enabled: &enabled,
		Tags:    o.Tags,
	}
}

type output struct {
	ID           *string                `json:"outputId"`
	CollectionID *string                `json:"collectionId"`
	Type         *string                `json:"type"`
	Config       map[string]interface{} `json:"config"`
	Enabled      *bool                  `json:"enabled"`
	Tags         map[string]string      `json:"tags,omitempty"`
}

func (o *output) toOutput() (Output, error) {
	switch *o.Type {
	case "webhook":
		return WebHookOutput{
			ID:                *o.ID,
			CollectionID:      *o.CollectionID,
			URL:               o.str("url"),
			BasicAuthUser:     o.str("basicAuthUser"),
			BasicAuthPass:     o.str("basicAuthPass"),
			CustomHeaderName:  o.str("customHeaderName"),
			CustomHeaderValue: o.str("customHeaderValue"),
			Disabled:          !*o.Enabled,
			Tags:              o.Tags,
		}, nil
	case "mqtt":
		return MQTTOutput{
			ID:               *o.ID,
			CollectionID:     *o.CollectionID,
			Endpoint:         o.str("endpoint"),
			DisableCertCheck: o.bool("disableCertCheck"),
			Username:         o.str("username"),
			Password:         o.str("password"),
			ClientID:         o.str("clientId"),
			TopicName:        o.str("topicName"),
			Disabled:         !*o.Enabled,
			Tags:             o.Tags,
		}, nil
	case "ifttt":
		return IFTTTOutput{
			ID:           *o.ID,
			CollectionID: *o.CollectionID,
			Key:          o.str("key"),
			EventName:    o.str("eventName"),
			AsIsPayload:  o.bool("asIsPayload"),
			Disabled:     !*o.Enabled,
			Tags:         o.Tags,
		}, nil
	case "udp":
		return UDPOutput{
			ID:           *o.ID,
			CollectionID: *o.CollectionID,
			Host:         o.str("host"),
			Port:         o.int("port"),
			Disabled:     !*o.Enabled,
			Tags:         o.Tags,
		}, nil
	}
	return nil, fmt.Errorf("unknown output type %q", *o.Type)
}

func (o *output) str(name string) string {
	s, _ := o.Config[name].(string)
	return s
}

func (o *output) bool(name string) bool {
	b, _ := o.Config[name].(bool)
	return b
}

func (o *output) int(name string) int {
	b, _ := o.Config[name].(float64)
	return int(b)
}
