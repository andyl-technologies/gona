package gona

// Location is an API response message of available deployment locations
type Location struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IATACode  string `json:"iata_code"`
	Continent string `json:"continent"`
	Flag      string `json:"flat"`
	Disabled  int    `json:"disabled"`
}

// GetLocations public method on Client to get a list of locations
func (c *Client) GetLocations() ([]Location, error) {
	r := make([]Location, 0)
	if err := c.get("cloud/locations", &r); err != nil {
		return nil, err
	}
	return r, nil
}
