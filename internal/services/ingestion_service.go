package services

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Abhaythakor/SigMap/internal/integrations/ipinfo"
	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type IngestionService struct {
	Repo         *repositories.DomainRepository
	IPInfoClient *ipinfo.Client
}

func NewIngestionService(repo *repositories.DomainRepository, ipInfoClient *ipinfo.Client) *IngestionService {
	return &IngestionService{Repo: repo, IPInfoClient: ipInfoClient}
}

type ScanResult struct {
	Domain     string `json:"domain"`
	URL        string `json:"url"`
	Technology string `json:"technology"`
	Confidence string `json:"confidence"`
	Source     string `json:"source"`
	Version    string `json:"version,omitempty"`
}

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

		// Infrastructure Enrichment on new/updated domain
		go s.LookupInfrastructure(ctx, domainID, res.Domain)

		confInt := 50
		switch strings.ToLower(res.Confidence) {
		case "high":
			confInt = 100
		case "medium":
			confInt = 70
		case "low":
			confInt = 40
		}

		err = s.Repo.AddDetection(ctx, domainID, res.Technology, res.URL, res.Version, confInt, res.Source)
		if err != nil {
			log.Printf("Error adding detection %s for %s: %v", res.Technology, res.Domain, err)
		}
	}

	return scanner.Err()
}

func (s *IngestionService) LookupInfrastructure(ctx context.Context, domainID int, domainName string) {
	// 1. Resolve IP
	ips, err := net.LookupIP(domainName)
	if err != nil || len(ips) == 0 {
		log.Printf("Infra: Could not resolve IP for %s", domainName)
		return
	}
	ip := ips[0].String()

	// 2. Fetch IP Details
	details, err := s.IPInfoClient.GetIPDetails(ctx, ip)
	if err != nil {
		log.Printf("Infra: IPInfo lookup failed for %s (%s): %v", domainName, ip, err)
		return
	}

	// 3. Parse ASN
	asn := 0
	asnOrg := details.Org
	if strings.HasPrefix(details.Org, "AS") {
		parts := strings.SplitN(details.Org, " ", 2)
		val, _ := strconv.Atoi(strings.TrimPrefix(parts[0], "AS"))
		asn = val
		if len(parts) > 1 {
			asnOrg = parts[1]
		}
	}

	// 4. Update DB
	_, err = s.Repo.Pool.Exec(ctx, `
		UPDATE domains SET 
			ip_address = $1, 
			cloud_provider = $2, 
			asn = $3, 
			asn_org = $4 
		WHERE id = $5
	`, ip, details.CloudProvider, asn, asnOrg, domainID)
	if err != nil {
		log.Printf("Infra: Failed to update DB for %s: %v", domainName, err)
	}
}
