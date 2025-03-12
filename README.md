# DeviceID

A Go package for generating and managing unique device identifiers. This package provides functionality to create, store, and verify persistent device IDs based on system-specific information.

## Features

- Cross-platform support (Windows, macOS, Linux)
- Persistent device ID storage
- SHA-256 based device ID generation
- Configurable storage location
- Thread-safe operations

## Installation

```bash
go get github.com/art3mis/deviceid
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/art3mis/deviceid"
)

func main() {
    // Create a new manager with default configuration
    manager := deviceid.NewManager(deviceid.Config{})
    
    // Get or generate a device ID
    id, err := manager.VerifyDeviceID()
    if err != nil {
        log.Fatalf("Failed to get device ID: %v", err)
    }
    
    fmt.Printf("Device ID: %s\n", id)
}
```

### Custom Configuration

You can customize the storage location and filename:

```go
config := deviceid.Config{
    StorageDir: "/custom/path",
    IDFileName: "custom-device-id",
}
manager := deviceid.NewManager(config)
```

## API Reference

### Types

#### `Config`
```go
type Config struct {
    StorageDir  string // Directory where the device ID file will be stored
    IDFileName  string // Name of the device ID file
}
```

#### `Manager`
```go
type Manager struct {}
```

### Functions

#### `NewManager(config Config) *Manager`
Creates a new device ID manager with the given configuration.

#### `(m *Manager) GenerateDeviceID() (string, error)`
Generates a new device ID based on system information.

#### `(m *Manager) SaveDeviceID(deviceID string) error`
Saves the device ID to the configured location.

#### `(m *Manager) VerifyDeviceID() (string, error)`
Verifies the existing device ID or generates a new one if needed.

#### `IsValidSHA256(s string) bool`
Checks if a string is a valid SHA-256 hash.

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.