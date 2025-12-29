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

# Create a new project (interactive)
amp projects create --org myorg

# Create a project with flags
amp projects create --org myorg --display-name "My Project" --pipeline default

# Delete a project
amp projects delete myproject --org myorg

# Skip confirmation with --force
amp projects delete myproject --org myorg --force

# Output as JSON
amp projects list --org myorg --output json
```

### Agents

```bash
# List agents in a project
amp agents list --org myorg --project myproject

# Get agent details
amp agents get myagent --org myorg --project myproject

# Create a new agent (interactive)
amp agents create --org myorg --project myproject

# Create an external agent with flags
amp agents create --display-name "My Agent" --provisioning external --org myorg --project myproject

# Create an internal agent with flags
amp agents create \
  --display-name "My Python Agent" \
  --provisioning internal \
  --repo-url https://github.com/user/repo \
  --branch main \
  --language python \
  --language-version 3.11 \
  --subtype chat-api \
  --org myorg --project myproject

# Delete an agent
amp agents delete myagent --org myorg --project myproject

# Skip confirmation with --force
amp agents delete myagent --org myorg --project myproject --force

# Output as JSON
amp agents get myagent --output json
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
| `amp projects get <name>` | Get details of a specific project |
| `amp projects create` | Create a new project |
| `amp projects delete <name>` | Delete a project |
| `amp agents list` | List all agents in a project |
| `amp agents get <name>` | Get details of a specific agent |
| `amp agents create` | Create a new agent |
| `amp agents delete <name>` | Delete an agent |
| `amp config list` | Show all configuration |
| `amp config set <key> <value>` | Set a configuration value |
| `amp config get <key>` | Get a configuration value |

## Flags

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--org` | `-o` | Organization name |
| `--project` | `-p` | Project name |
| `--output` | | Output format: `table` or `json` |

### Command-Specific Flags

#### `projects create`
| Flag | Description |
|------|-------------|
| `--name` | Project name (auto-generated if not provided) |
| `--display-name` | Display name for the project |
| `--description` | Project description |
| `--pipeline` | Deployment pipeline name |

#### `agents create`
| Flag | Description |
|------|-------------|
| `--name` | Agent name (auto-generated if not provided) |
| `--display-name` | Display name for the agent |
| `--description` | Agent description |
| `--provisioning` | Provisioning type: `internal` or `external` |
| `--repo-url` | Repository URL (for internal agents) |
| `--branch` | Git branch (default: main) |
| `--app-path` | App path in repository (default: /) |
| `--subtype` | Agent subtype: `chat-api` or `custom-api` |
| `--language` | Programming language (python, nodejs, java, go, ballerina) |
| `--language-version` | Language version |

#### Delete Commands
| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Skip confirmation prompts |

## Interactive Mode

Running `amp` without arguments launches an interactive shell:

```bash
$ amp
╭─────────────────────────────────────────────────────────╮
│                    WSO2 AMP CLI                         │
│          AI Agent Management Platform                   │
╰─────────────────────────────────────────────────────────╯

amp> orgs list
amp> projects list
amp> exit
```

## Development

```bash
# Run without building
go run ./cmd/amp agents list --org test --project test

# Build
go build -o amp ./cmd/amp

# Run tests
go test ./...

# Vet
go vet ./...
```

## License

Apache 2.0
