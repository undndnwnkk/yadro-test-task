package dns

import (
	"os"
	"testing"
)

func TestDNSGet(t *testing.T) {
	tempFile := "test_resolv.conf"
	os.WriteFile(tempFile, []byte("nameserver 8.8.8.8\n"), 0644)
	defer os.Remove(tempFile)

	service := NewService(tempFile)

	servers, _ := service.GetServers()
	if len(servers) != 1 || servers[0] != "8.8.8.8" {
		t.Errorf("expected 8.8.8.8, got %v", servers)
	}
}

func TestDNSAdd(t *testing.T) {
	tempFile := "test_resolv.conf"
	os.WriteFile(tempFile, []byte("nameserver 8.8.8.8\n"), 0644)
	defer os.Remove(tempFile)

	service := NewService(tempFile)

	service.AddServer("1.1.1.1")
	servers, _ := service.GetServers()
	if len(servers) != 2 {
		t.Error("failed to add server")
	}
}

func TestDNSRemove(t *testing.T) {
	tempFile := "test_resolv.conf"
	os.WriteFile(tempFile, []byte("nameserver 8.8.8.8\n"), 0644)
	defer os.Remove(tempFile)

	service := NewService(tempFile)

	service.RemoveServer("8.8.8.8")
	servers, _ := service.GetServers()
	if len(servers) != 0 {
		t.Error("failed to remove server")
	}
}
