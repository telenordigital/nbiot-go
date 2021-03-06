package nbiot

import (
	"fmt"
	"time"
)

// Collection represents a collection.
type Collection struct {
	ID        string            `json:"collectionId"`
	TeamID    string            `json:"teamId,omitempty"`
	FieldMask *FieldMask        `json:"fieldMask,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
}

// FieldMask indicates which fields will be masked from API responses.
type FieldMask struct {
	IMSI     *bool `json:"imsi"`
	IMEI     *bool `json:"imei"`
	Location *bool `json:"location"`
	MSISDN   *bool `json:"msisdn"`
}

// Collection gets a collection.
func (c *Client) Collection(id string) (Collection, error) {
	var collection Collection
	err := c.get("/collections/"+id, &collection)
	return collection, err
}

// Collections gets all collections that the user has access to.
func (c *Client) Collections() ([]Collection, error) {
	var collections struct {
		Collections []Collection `json:"collections"`
	}
	err := c.get("/collections", &collections)
	return collections.Collections, err
}

// CreateCollection creates a collection.
func (c *Client) CreateCollection(collection Collection) (Collection, error) {
	err := c.create("/collections", &collection)
	return collection, err
}

// UpdateCollection updates a collection.
// No tags are deleted, only added or updated.
func (c *Client) UpdateCollection(collection Collection) (Collection, error) {
	err := c.update("/collections/"+collection.ID, &collection)
	return collection, err
}

// DeleteCollectionTag deletes a tag from a collection.
func (c *Client) DeleteCollectionTag(id, name string) error {
	return c.delete(fmt.Sprintf("/collections/%s/tags/%s", id, name))
}

// DeleteCollection deletes a collection.
func (c *Client) DeleteCollection(id string) error {
	return c.delete("/collections/" + id)
}

// CollectionData returns all the stored data for the collection.
func (c *Client) CollectionData(collectionID string, since time.Time, until time.Time, limit int) ([]OutputDataMessage, error) {
	var s, u int64
	if !since.IsZero() {
		s = since.UnixNano() / int64(time.Millisecond)
	}
	if !until.IsZero() {
		u = until.UnixNano() / int64(time.Millisecond)
	}

	var data struct {
		Messages []OutputDataMessage `json:"messages"`
	}
	err := c.get(fmt.Sprintf("/collections/%s/data?since=%d&until=%d&limit=%d", collectionID, s, u, limit), &data)
	return data.Messages, err
}
