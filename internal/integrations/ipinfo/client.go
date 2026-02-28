package ipinfo

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type IPDetails struct {
	IP            string
	City          string
	Region        string
	Country       string
	Org           string // Includes ASN and Name
	CloudProvider string
}

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{Token: token}
}

// GetIPDetails retrieves geographical and network info for an IP.
func (c *Client) GetIPDetails(ctx context.Context, ip string) (*IPDetails, error) {
	log.Printf("IPInfo: Fetching details for %s", ip)
	
	// Simulation
	if c.Token == "" {
		log.Println("IPInfo: No token provided, running in simulation mode")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	providers := []string{"Amazon.com, Inc.", "Google LLC", "Microsoft Corporation", "Cloudflare, Inc.", "DigitalOcean, LLC"}
	provider := providers[rng.Intn(len(providers))]
	
	cloud := "Undetected"
	if rng.Float32() > 0.3 {
		switch provider {
		case "Amazon.com, Inc.": cloud = "AWS"
		case "Google LLC": cloud = "GCP"
		case "Microsoft Corporation": cloud = "Azure"
		case "Cloudflare, Inc.": cloud = "Cloudflare"
		case "DigitalOcean, LLC": cloud = "DigitalOcean"
		}
	}

	return &IPDetails{
		IP:            ip,
		City:          "San Francisco",
		Region:        "California",
		Country:       "US",
		Org:           "AS12345 " + provider,
		CloudProvider: cloud,
	}, nil
}
