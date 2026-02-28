package services

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type IngestionService struct {
	Repo *repositories.DomainRepository
}

func NewIngestionService(repo *repositories.DomainRepository) *IngestionService {
	return &IngestionService{Repo: repo}
}

// ScanResult matches the JSON format in testDir
type ScanResult struct {
	Domain     string `json:"domain"`
	URL        string `json:"url"`
	Technology string `json:"technology"`
	Confidence string `json:"confidence"`
	Source     string `json:"source"`
	Version    string `json:"version,omitempty"`
}

// IngestFromDirectory reads all .json files in a directory and saves them.
func (s *IngestionService) IngestFromDirectory(ctx context.Context, dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(dirPath, file.Name())
			log.Printf("Ingesting file: %s", filePath)
			if err := s.ingestFile(ctx, filePath); err != nil {
				log.Printf("Error ingesting %s: %v", filePath, err)
			}
		}
	}
	return nil
}

func (s *IngestionService) ingestFile(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var res ScanResult
		line := scanner.Text()
		if line == "" {
			continue
		}

		if err := json.Unmarshal([]byte(line), &res); err != nil {
			log.Printf("Skip invalid JSON line in %s: %v", filePath, err)
			continue
		}

		domainID, err := s.Repo.EnsureDomain(ctx, res.Domain)
		if err != nil {
			return err
		}

		// Map string confidence to int
		confInt := 50
		switch strings.ToLower(res.Confidence) {
		case "high":
			confInt = 100
		case "medium":
			confInt = 70
		case "low":
			confInt = 40
		}

		// Add detection (repo handles "Tech:Version" splitting automatically)
		err = s.Repo.AddDetection(ctx, domainID, res.Technology, res.URL, res.Version, confInt, res.Source)
		if err != nil {
			log.Printf("Error adding detection %s for %s: %v", res.Technology, res.Domain, err)
		}
	}

	return scanner.Err()
}

// IngestSampleData remains for mock testing if needed
func (s *IngestionService) IngestSampleData(ctx context.Context) error {
    // Logic kept for backup/reference
    return nil
}

func (s *IngestionService) LookupInfrastructure(ctx context.Context, domainID int, domainName string) {
	// Simple mock infrastructure logic
	ips := []string{"34.212.12.1", "52.4.15.22", "104.16.24.5", "13.248.155.12"}
	idx := len(domainName) % len(ips)
	s.Repo.Pool.Exec(ctx, "UPDATE domains SET ip_address = $1, cloud_provider = 'Detected' WHERE id = $2", ips[idx], domainID)
}
