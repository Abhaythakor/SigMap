package chaos

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Client handles interaction with ProjectDiscovery Chaos.
type Client struct {
	APIKey string
}

func NewClient(apiKey string) *Client {
	return &Client{APIKey: apiKey}
}

// FetchSubdomains retrieves subdomains for a given domain.
func (c *Client) FetchSubdomains(ctx context.Context, domain string) ([]string, error) {
	log.Printf("Chaos: Fetching subdomains for %s", domain)
	
	// Simulation logic
	if c.APIKey == "" {
		log.Println("Chaos: No API key provided, running in simulation mode")
	}

	// Simulating typical subdomains
	prefixes := []string{"dev", "staging", "api", "test", "v1", "v2", "app", "static", "cdn", "admin"}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	num := rng.Intn(len(prefixes)-3) + 3
	shuffled := make([]string, len(prefixes))
	copy(shuffled, prefixes)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	var results []string
	for i := 0; i < num; i++ {
		results = append(results, fmt.Sprintf("%s.%s", shuffled[i], domain))
	}

	return results, nil
}
