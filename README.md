# EmBackUp

A robust multi-database backup utility written in Go.

## Overview

EmBackUp is a command-line tool designed to simplify and automate backups across multiple database systems. Whether you're managing MySQL, PostgreSQL, MongoDB, or other databases, EmBackUp provides a unified interface for scheduling, executing, and managing database backups.

## Features

- **Multi-Database Support** - Back up multiple database types from a single application
- **Flexible Configuration** - Simple configuration file-based setup
- **Command-Line Interface** - Easy-to-use CLI for manual backups and management
- **GUI Support** - Optional graphical interface for users who prefer visual tools
- **Automated Scheduling** - Schedule regular backups with cron-like functionality
- **Lightweight & Fast** - Written in Go for minimal resource usage and fast execution

## Project Structure

```
.
├── cmd/                    # Command-line application entry points
├── config/                 # Configuration handling and defaults
├── gui/                    # Graphical user interface components
├── internal/               # Internal packages and utilities
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── LICENSE                 # MIT License
```

## Installation

### Prerequisites

- Go 1.19 or higher

### Building from Source

```bash
git clone https://github.com/EmilioQ74/EmBackUp.git
cd EmBackUp
go build -o embackup
```

## Usage

### Command Line

```bash
./embackup [command] [options]
```

For detailed command documentation, run:

```bash
./embackup --help
```

### Configuration

Configuration files should be placed in the `config/` directory. See the configuration section of the documentation for detailed setup instructions.

## Development

### Project Layout

- **cmd/** - Command definitions and CLI logic
- **config/** - Configuration file parsing and validation
- **gui/** - UI components (if using the graphical interface)
- **internal/** - Core backup logic, database drivers, and utilities

### Building

```bash
go build -o embackup
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open issues for bug reports and feature requests, or submit pull requests with improvements.

## Support

For issues, questions, or feature requests, please open an [issue](https://github.com/EmilioQ74/EmBackUp/issues) on GitHub.

---

**Author:** Emilio  
**Repository:** [EmilioQ74/EmBackUp](https://github.com/EmilioQ74/EmBackUp)
