package cat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kontentski/develops-today-task/config"
	"github.com/Kontentski/develops-today-task/pkg/logging"
)

// catAPI implements the cat API client
type catAPI struct {
	http   *http.Client
	logger logging.Logger
	cfg    *config.Config
}

// Options is used to parameterize catAPI using New
type Options struct {
	Logger logging.Logger
	Config *config.Config
}

// Breed represents a cat breed from TheCatAPI
type Breed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"`
	Temperament string `json:"temperament"`
}

// New creates a new catAPI instance
func New(options *Options) *catAPI {
	return &catAPI{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: options.Logger.Named("CatAPI"),
		cfg:    options.Config,
	}
}

// GetBreeds fetches all cat breeds from TheCatAPI
func (c *catAPI) GetBreeds() ([]Breed, error) {
	c.logger.Debug("fetching cat breeds")

	// get with context
	resp, err := c.http.Get(fmt.Sprintf("%s/breeds", c.cfg.CatAPI.URL))
	if err != nil {
		c.logger.Error("failed to fetch breeds", "err", err)
		return nil, fmt.Errorf("failed to fetch breeds: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("thecatapi responded with non-200 status", "status", resp.Status)
		return nil, fmt.Errorf("thecatapi responded with status: %s", resp.Status)
	}

	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		c.logger.Error("failed to decode breeds response", "err", err)
		return nil, fmt.Errorf("failed to decode breeds response: %w", err)
	}

	c.logger.Info("successfully fetched cat breeds", "count", len(breeds))
	return breeds, nil
}
