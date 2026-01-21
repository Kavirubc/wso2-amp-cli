# amp-cli

A command-line tool for managing the WSO2 AI Agent Management Platform.

## Installation

### From Source (requires Go 1.21+)

```bash
git clone https://github.com/Kavirubc/wso2-amp-cli.git
cd wso2-amp-cli
go build -o amp ./cmd/amp
```

### Using go install

```bash
go install github.com/Kavirubc/wso2-amp-cli/cmd/amp@latest
```

### Pre-built Binaries

Download from the [Releases](https://github.com/Kavirubc/wso2-amp-cli/releases) page.

## Quick Start

```bash
# 1. Run the setup wizard
amp login

# 2. List your agents
amp agents list

# 3. View agent details
amp agents get my-agent
```

## Configuration

### Using Login Command (Recommended)

```bash
amp login
```

This interactive wizard will guide you through:
- Setting the API server URL
- Authenticating with your token
- Selecting default organization and project

### Manual Configuration

```bash
# Set API server URL
amp config set api_url https://your-server.com/api/v1

# Set authentication
amp config set api_key_header Authorization
amp config set api_key "Bearer your-token"

# Set defaults for convenience
amp config set default_org your-org
amp config set default_project your-project

# View current configuration
amp config show
```

### Configuration File

Settings are stored in `~/.amp/config.yaml`

| Key | Description |
|-----|-------------|
| `api_url` | API server URL |
| `api_key_header` | Authentication header name |
| `api_key` | Authentication token |
| `default_org` | Default organization |
| `default_project` | Default project |

## Commands

### Authentication

#### `amp login`
Interactive setup wizard for initial configuration.

```bash
amp login

# Non-interactive mode
amp login --api-url https://api.example.com --token your-token
```

#### `amp logout`
Clear stored credentials.

```bash
amp logout

# Skip confirmation
amp logout --force
```

### Organizations

#### `amp orgs list`
List all organizations.

```bash
amp orgs list
amp orgs list --output json
```

#### `amp orgs get <name>`
Get organization details.

```bash
amp orgs get my-org
```

### Projects

#### `amp projects list`
List all projects in an organization.

```bash
amp projects list
amp projects list --org my-org
```

#### `amp projects create`
Create a new project.

```bash
# Interactive mode
amp projects create

# With flags
amp projects create \
  --display-name "My Project" \
  --description "Project description" \
  --pipeline default
```

#### `amp projects delete <name>`
Delete a project.

```bash
amp projects delete my-project

# Skip confirmation
amp projects delete my-project --force
```

### Agents

#### `amp agents list`
List all agents in a project.

```bash
amp agents list
amp agents list --org my-org --project my-project
```

#### `amp agents get <name>`
Get agent details.

```bash
amp agents get my-agent
```

#### `amp agents create`
Create a new agent.

```bash
# Interactive mode
amp agents create

# External agent
amp agents create \
  --display-name "My Agent" \
  --provisioning external

# Internal agent with repository
amp agents create \
  --display-name "My Python Agent" \
  --provisioning internal \
  --repo-url https://github.com/user/repo \
  --branch main \
  --language python \
  --language-version 3.11 \
  --subtype chat-api
```

#### `amp agents delete <name>`
Delete an agent.

```bash
amp agents delete my-agent

# Skip confirmation
amp agents delete my-agent --force
```

#### `amp agents token`
Generate a JWT token for an agent.

```bash
amp agents token --agent my-agent

# Custom expiration
amp agents token --agent my-agent --expires-in 24h
```

#### `amp agents logs`
View runtime logs for a deployed agent.

```bash
amp agents logs --agent my-agent --env development

# Filter by log level
amp agents logs --agent my-agent --env dev --level ERROR,WARN

# Search logs
amp agents logs --agent my-agent --env dev --search "error"

# Limit results
amp agents logs --agent my-agent --env dev --limit 50
```

### Builds

#### `amp builds list`
List all builds for an agent.

```bash
amp builds list --agent my-agent
```

#### `amp builds get <name>`
Get build details with steps.

```bash
amp builds get build-123 --agent my-agent
```

#### `amp builds trigger`
Trigger a new build.

```bash
amp builds trigger --agent my-agent

# Trigger with specific commit
amp builds trigger --agent my-agent --commit abc123
```

#### `amp builds logs <name>`
View build logs.

```bash
amp builds logs build-123 --agent my-agent
```

### Deployments

#### `amp deploy`
Deploy an agent to an environment.

```bash
amp deploy --agent my-agent --image build-123

# With environment variables
amp deploy --agent my-agent --image build-123 \
  --set-env API_KEY=secret \
  --set-env DEBUG=true
```

#### `amp deployments list`
List all deployments for an agent.

```bash
amp deployments list --agent my-agent
```

#### `amp deployments endpoints`
List endpoints for a deployed agent.

```bash
amp deployments endpoints --agent my-agent

# Filter by environment
amp deployments endpoints --agent my-agent --env production
```

### Configuration

#### `amp config show`
Display all configuration settings.

```bash
amp config show
```

#### `amp config set <key> <value>`
Set a configuration value.

```bash
amp config set api_url https://api.example.com
```

#### `amp config get <key>`
Get a specific configuration value.

```bash
amp config get api_url
```

#### `amp config reset`
Reset configuration to defaults.

```bash
amp config reset
```

### Other Commands

#### `amp version`
Display version information.

```bash
amp version
```

#### `amp completion`
Generate shell completion scripts.

```bash
# Bash
amp completion bash > /etc/bash_completion.d/amp

# Zsh
amp completion zsh > "${fpath[1]}/_amp"

# Fish
amp completion fish > ~/.config/fish/completions/amp.fish
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--org` | `-o` | Organization name |
| `--project` | `-p` | Project name |
| `--output` | | Output format: `table` or `json` |
| `--verbose` | `-v` | Enable debug output |
| `--help` | `-h` | Show help |

## Interactive Mode

Run `amp` without arguments to start an interactive shell:

```
$ amp
╭─────────────────────────────────────────────────────────╮
│                    WSO2 AMP CLI                         │
│          AI Agent Management Platform                   │
╰─────────────────────────────────────────────────────────╯

amp> orgs list
amp> agents list
amp> exit
```

## Examples

### Deploy a New Agent

```bash
# 1. Create the agent
amp agents create \
  --display-name "My API Agent" \
  --provisioning internal \
  --repo-url https://github.com/user/my-agent \
  --branch main \
  --language python \
  --language-version 3.11 \
  --subtype custom-api

# 2. Wait for build to complete
amp builds list --agent my-api-agent

# 3. Deploy to development
amp deploy --agent my-api-agent --image build-xyz

# 4. Check deployment status
amp deployments list --agent my-api-agent

# 5. Get the endpoint URL
amp deployments endpoints --agent my-api-agent --env development
```

### Monitor an Agent

```bash
# View recent builds
amp builds list --agent my-agent

# Check build details
amp builds get build-123 --agent my-agent

# View build logs
amp builds logs build-123 --agent my-agent

# View runtime logs
amp agents logs --agent my-agent --env development --since 1h
```

### Generate Agent Token

```bash
# Generate token for API calls
amp agents token --agent my-agent --expires-in 7d

# Use the token in your application
export AGENT_TOKEN=$(amp agents token --agent my-agent --output json | jq -r .token)
```

## Troubleshooting

### Authentication Issues

```
✗ Authentication failed
```

**Solution:** Run `amp login` to reconfigure your credentials, or check:
```bash
amp config show
```

### Connection Errors

```
✗ Cannot connect to API server
```

**Solution:** Verify the API URL is correct:
```bash
amp config get api_url
```

### Resource Not Found

```
✗ Agent 'my-agent' not found
```

**Solution:** List available resources:
```bash
amp agents list
amp projects list
amp orgs list
```

### Debug Mode

Enable verbose output to see detailed request/response information:
```bash
amp agents list --verbose
```

## Development

```bash
# Run without building
go run ./cmd/amp agents list

# Build
go build -o amp ./cmd/amp

# Run tests
go test ./...

# Vet
go vet ./...
```

## License

Apache 2.0
