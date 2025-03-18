// Package deviceid provides functionality for generating and managing unique device identifiers.
// It creates and maintains persistent device IDs based on system-specific information,
// ensuring consistent identification across application restarts.
package deviceid

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// DefaultIDFileName is the name of the file where the device ID is stored
const DefaultIDFileName = ".device_id"

// Config holds the configuration for device ID management
type Config struct {
	// Directory where the device ID file will be stored
	// If empty, defaults to ~/.parity
	StorageDir string
	// Name of the device ID file
	// If empty, defaults to .device_id
	IDFileName string
}

// Manager handles device ID operations
type Manager struct {
	config Config
}

// NewManager creates a new device ID manager with the given configuration
func NewManager(config Config) *Manager {
	if config.IDFileName == "" {
		config.IDFileName = DefaultIDFileName
	}
	return &Manager{config: config}
}

// getSystemInfo retrieves system-specific information based on the operating system
func getSystemInfo() (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("wmic", "csproduct", "get", "UUID")
	case "darwin":
		cmd = exec.Command("ioreg", "-d2", "-c", "IOPlatformExpertDevice")
	default: // Linux
		cmd = exec.Command("cat", "/etc/machine-id")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get system info: %w", err)
	}
	return string(output), nil
}

// GenerateDeviceID creates a new device ID based on system information
func (m *Manager) GenerateDeviceID() (string, error) {
	info, err := getSystemInfo()
	if err != nil {
		return "", fmt.Errorf("failed to generate device ID: %w", err)
	}

	hash := sha256.Sum256([]byte(info))
	return hex.EncodeToString(hash[:]), nil
}

// GetDeviceIDPath returns the full path where the device ID file should be stored
func (m *Manager) GetDeviceIDPath() (string, error) {
	var basePath string

	if m.config.StorageDir != "" {
		basePath = m.config.StorageDir
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		basePath = filepath.Join(home, ".parity")
	}

	return filepath.Join(basePath, m.config.IDFileName), nil
}

// SaveDeviceID stores the device ID in the configured location
func (m *Manager) SaveDeviceID(deviceID string) error {
	if !IsValidSHA256(deviceID) {
		return fmt.Errorf("invalid device ID format")
	}

	path, err := m.GetDeviceIDPath()
	if err != nil {
		return fmt.Errorf("failed to get device ID path: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(deviceID), 0o600); err != nil {
		return fmt.Errorf("failed to write device ID: %w", err)
	}

	return nil
}

// VerifyDeviceID checks for an existing device ID and generates a new one if needed
func (m *Manager) VerifyDeviceID() (string, error) {
	path, err := m.GetDeviceIDPath()
	if err != nil {
		return "", fmt.Errorf("failed to get device ID path: %w", err)
	}

	// Check if device ID exists
	deviceID, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Generate new device ID if it doesn't exist
			newID, err := m.GenerateDeviceID()
			if err != nil {
				return "", fmt.Errorf("failed to generate new device ID: %w", err)
			}
			if err := m.SaveDeviceID(newID); err != nil {
				return "", fmt.Errorf("failed to save new device ID: %w", err)
			}
			return newID, nil
		}
		return "", fmt.Errorf("failed to read device ID: %w", err)
	}

	// Validate the stored device ID format
	storedID := string(deviceID)
	if !IsValidSHA256(storedID) {
		// If invalid format, generate a new one
		newID, err := m.GenerateDeviceID()
		if err != nil {
			return "", fmt.Errorf("failed to generate new device ID: %w", err)
		}
		if err := m.SaveDeviceID(newID); err != nil {
			return "", fmt.Errorf("failed to save new device ID: %w", err)
		}
		return newID, nil
	}

	return storedID, nil
}

// IsValidSHA256 checks if a string is a valid SHA256 hash
func IsValidSHA256(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}
