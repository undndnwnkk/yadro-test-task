package dns

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Service struct {
	ConfigPath string
}

func NewService(configPath string) *Service {
	return &Service{ConfigPath: configPath}
}

func (s *Service) GetServers() ([]string, error) {
	log.Printf("DEBUG: Reading DNS servers from %s", s.ConfigPath)

	file, err := os.Open(s.ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "nameserver") {
			parts := strings.Fields(line)

			if len(parts) >= 2 {
				servers = append(servers, parts[1])
			}
		}
	}

	return servers, scanner.Err()
}

func (s *Service) AddServer(ip string) error {
	log.Printf("DEBUG: Attempting to add DNS server: %s", ip)
	servers, err := s.GetServers()
	if err != nil {
		return err
	}

	for _, existing := range servers {
		if existing == ip {
			return fmt.Errorf("server %s already exists", ip)
		}
	}

	file, err := os.OpenFile(s.ConfigPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("nameserver %s\n", ip))
	return err
}

func (s *Service) RemoveServer(ip string) error {
	log.Printf("DEBUG: Attempting to remove DNS server: %s", ip)
	file, err := os.Open(s.ConfigPath)
	if err != nil {
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	isFound := false

	for scanner.Scan() {
		currentLine := scanner.Text()

		if strings.TrimSpace(line) == "nameserver " + ip {
			isFound = true
			continue
		}

		lines = append(lines, currentLine)
	}
	file.Close()

	if !isFound {
		return fmt.Errorf("server %s not found", ip)
	}

	return os.WriteFile(s.ConfigPath, []byte(strings.Join(lines, "\n") + "\n"), 0644)
}

