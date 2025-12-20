# amp-cli

A CLI tool for managing the WSO2 AI Agent Management Platform.

## Installation

### From Source (requires Go 1.21+)

```bash
# Clone the repository
git clone https://github.com/Kavirubc/wso2-amp-cli.git
cd wso2-amp-cli

# Install globally (produces 'amp' binary)
go install ./cmd/amp

# Or build locally
go build -o amp ./cmd/amp
```

### Using go install

```bash
go install github.com/Kavirubc/wso2-amp-cli/cmd/amp@latest
```

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/Kavirubc/wso2-amp-cli/releases) page.

## Configuration

Configure the CLI before first use:

```bash
# Set your API server URL
amp config set api_url http://your-server:8080

# Set your API key (if required)
amp config set api_key your-api-key

# Set defaults for convenience
amp config set default_org your-org-name
amp config set default_project your-project-name

# View current configuration
amp config list
```

Configuration is stored in `~/.amp/config.yaml`

## Usage

### Organizations

```bash
# List all organizations
amp orgs list

# Output as JSON
amp orgs list --output json
```

### Projects

```bash
# List projects in an organization
amp projects list --org myorg

# Get project details
amp projects get myproject --org myorg

# Delete a project
amp projects delete myproject --org myorg

# Skip confirmation with --force
amp projects delete myproject --org myorg --force

# Output as JSON
amp projects list --org myorg --output json
```

### Agents

```bash
# Using flags
amp agents list --org myorg --project myproject

# Using defaults (if configured)
amp agents list

# Output as JSON
amp agents list --output json
```

### Configuration

```bash
amp config list              # List all settings
amp config get api_url       # Get a specific setting
amp config set api_url URL   # Set a setting
```

## Available Commands

| Command | Description |
|---------|-------------|
| `amp orgs list` | List all organizations |
| `amp projects list` | List all projects in an organization |
| `amp projects get` | Get details of a specific project |
| `amp projects delete` | Delete a project |
| `amp agents list` | List all agents in a project |
| `amp config list` | Show all configuration |
| `amp config set` | Set a configuration value |
| `amp config get` | Get a configuration value |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--org` | `-o` | Organization name |
| `--project` | `-p` | Project name |
| `--output` | | Output format: `table` or `json` |
| `--force` | `-f` | Skip confirmation prompts |

## Development

```bash
# Run without building
go run ./cmd/amp agents list --org test --project test

# Build
go build -o amp ./cmd/amp

# Run tests
go test ./...
```

## License

Apache 2.0
