package horde

import "fmt"

// Output represents a data output for a collection
type Output struct {
	OutputID     *string                `json:"outputId"`
	CollectionID *string                `json:"collectionId"`
	Type         *string                `json:"type"`
	Config       map[string]interface{} `json:"config"`
}

// Output retrieves an output
func (c *Client) Output(collectionID, outputID string) (Output, error) {
	var output Output
	err := c.get(fmt.Sprintf("/collections/%s/outputs/%s", collectionID, outputID), &output)
	return output, err
}

// Outputs retrieves a list of outputs on a collection
func (c *Client) Outputs(collectionID string) ([]Output, error) {
	var outputs struct {
		Outputs []Output `json:"outputs"`
	}
	err := c.get(fmt.Sprintf("/collections/%s/outputs", collectionID), &outputs)
	return outputs.Outputs, err
}

// CreateOutput creates an output
func (c *Client) CreateOutput(collectionID string, output Output) (Output, error) {
	err := c.create(fmt.Sprintf("/collections/%s/outputs", collectionID), &output)
	return output, err
}

// UpdateOutput updates an output. The type field can't be modified
func (c *Client) UpdateOutput(collectionID string, output Output) (Output, error) {
	err := c.update(fmt.Sprintf("/collections/%s/outputs/%s", collectionID, *output.OutputID), &output)
	return output, err
}

// DeleteOutput removes an output
func (c *Client) DeleteOutput(collectionID, outputID string) error {
	return c.delete(fmt.Sprintf("/collections/%s/outputs/%s", collectionID, outputID))
}
