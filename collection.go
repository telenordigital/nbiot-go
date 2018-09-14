package horde

// Collection represents a collection.
type Collection struct {
	CollectionID string            `json:"collectionId"`
	TeamID       *string           `json:"teamId"`
	Tags         map[string]string `json:"tags,omitempty"`
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
func (c *Client) UpdateCollection(collection Collection) (Collection, error) {
	err := c.update("/collections/"+collection.CollectionID, &collection)
	return collection, err
}

// DeleteCollection deletes a collection.
func (c *Client) DeleteCollection(id string) error {
	return c.delete("/collections/" + id)
}
