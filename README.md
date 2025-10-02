# VCX

Personal version control.

## Architecture

VCX is broadly divided into:
- **agent** - API Server, FSMonitor, DB
- **clients** - CLI, React web interface
- **tauri** - application wrapper for distribution

## Development Setup

### Prerequisites
- Go 1.19+
- Node.js 16+
- Air (for auto-reload): `go install github.com/cosmtrek/air@latest`
