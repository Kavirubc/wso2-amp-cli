# AMP CLI Command Reference

Complete command reference for the WSO2 AI Agent Management Platform CLI.

## Command Overview

| Command | Description | API Endpoint |
|---------|-------------|--------------|
| `amp login` | Authenticate and configure CLI | - |
| `amp logout` | Clear credentials | - |
| `amp version` | Print version info | - |
| `amp orgs list` | List organizations | `GET /orgs` |
| `amp orgs get` | Get organization details | `GET /orgs/{name}` |
| `amp orgs create` | Create organization | `POST /orgs` |
| `amp projects list` | List projects | `GET /orgs/{org}/projects` |
| `amp projects get` | Get project details | `GET /orgs/{org}/projects/{name}` |
| `amp projects create` | Create project | `POST /orgs/{org}/projects` |
| `amp projects delete` | Delete project | `DELETE /orgs/{org}/projects/{name}` |
| `amp projects pipeline` | Get deployment pipeline | `GET /orgs/{org}/projects/{name}/deployment-pipeline` |
| `amp agents list` | List agents | `GET /orgs/{org}/projects/{proj}/agents` |
| `amp agents get` | Get agent details | `GET /orgs/{org}/projects/{proj}/agents/{name}` |
| `amp agents create` | Create agent | `POST /orgs/{org}/projects/{proj}/agents` |
| `amp agents delete` | Delete agent | `DELETE /orgs/{org}/projects/{proj}/agents/{name}` |
| `amp agents token` | Generate JWT token | `POST .../agents/{name}/token` |
| `amp agents logs` | View runtime logs | `POST .../agents/{name}/runtime-logs` |
| `amp builds list` | List builds | `GET .../agents/{agent}/builds` |
| `amp builds get` | Get build details | `GET .../agents/{agent}/builds/{name}` |
| `amp builds trigger` | Trigger build | `POST .../agents/{agent}/builds` |
| `amp builds logs` | View build logs | `GET .../agents/{agent}/builds/{name}/build-logs` |
| `amp deployments list` | List deployments | `GET .../agents/{agent}/deployments` |
| `amp deployments endpoints` | List endpoints | `GET .../agents/{agent}/endpoints` |
| `amp deploy` | Deploy agent | `POST .../agents/{agent}/deployments` |
| `amp environments list` | List environments | `GET /orgs/{org}/environments` |
| `amp dataplanes list` | List data planes | `GET /orgs/{org}/data-planes` |
| `amp pipelines list` | List pipelines | `GET /orgs/{org}/deployment-pipelines` |
| `amp pipelines get` | Get pipeline details | `GET /orgs/{org}/deployment-pipelines/{name}` |
| `amp config set` | Set config value | - |
| `amp config get` | Get config value | - |
| `amp config list` | List config | - |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--org` | `-o` | Organization name |
| `--project` | `-p` | Project name |
| `--output` | | Output format: `table` or `json` |
| `--verbose` | `-v` | Enable verbose output |
| `--help` | `-h` | Show help |

## Authentication

### Login

```bash
amp login
```

Interactive setup prompts for:
- API server URL
- Authentication token
- Default organization
- Default project

### Logout

```bash
amp logout
amp logout --force  # Skip confirmation
```

## Organizations

### List Organizations

```bash
amp orgs list
amp orgs list --output json
```

### Get Organization

```bash
amp orgs get my-org
amp orgs get my-org --output json
```

### Create Organization

```bash
amp orgs create my-org
amp orgs create --name my-org
amp orgs create my-org --output json
```

## Projects

### List Projects

```bash
amp projects list
amp projects list --org my-org
```

### Get Project

```bash
amp projects get my-project
amp projects get my-project --org my-org
```

### Create Project

```bash
amp projects create \
  --name my-project \
  --display-name "My Project" \
  --pipeline default-pipeline
```

### Delete Project

```bash
amp projects delete my-project
amp projects delete my-project --force
```

### Get Project Pipeline

```bash
amp projects pipeline my-project
```

## Agents

### List Agents

```bash
amp agents list
amp agents list --project my-project
```

### Get Agent

```bash
amp agents get my-agent
```

### Create Agent

```bash
# Interactive mode
amp agents create

# With flags
amp agents create \
  --name my-agent \
  --display-name "My Agent" \
  --provisioning external \
  --repo-url https://github.com/user/repo \
  --branch main \
  --subtype chat-api \
  --language python \
  --language-version "3.11"
```

### Delete Agent

```bash
amp agents delete my-agent
amp agents delete my-agent --force
```

### Generate Agent Token

```bash
amp agents token --agent my-agent
amp agents token --agent my-agent --expires-in 7d
```

### View Runtime Logs

```bash
amp agents logs --agent my-agent --env development
amp agents logs --agent my-agent --env production --since 1h --level ERROR
amp agents logs --agent my-agent --env development --search "error" --limit 50
```

## Builds

### List Builds

```bash
amp builds list --agent my-agent
```

### Get Build Details

```bash
amp builds get build-001 --agent my-agent
```

### Trigger Build

```bash
amp builds trigger --agent my-agent
amp builds trigger --agent my-agent --commit abc123
```

### View Build Logs

```bash
amp builds logs build-001 --agent my-agent
```

## Deployments

### List Deployments

```bash
amp deployments list --agent my-agent
```

### List Endpoints

```bash
amp deployments endpoints --agent my-agent
amp deployments endpoints --agent my-agent --env production
```

### Deploy Agent

```bash
amp deploy --agent my-agent --image sha256:abc123
amp deploy --agent my-agent --image sha256:abc123 --set-env API_KEY=xxx --set-env DEBUG=true
```

## Environments

```bash
amp environments list
amp envs list  # alias
```

## Data Planes

```bash
amp dataplanes list
amp dp list  # alias
```

## Pipelines

### List Pipelines

```bash
amp pipelines list
```

### Get Pipeline

```bash
amp pipelines get default-pipeline
```

## Configuration

### Set Value

```bash
amp config set api_url https://api.example.com
amp config set default_org my-org
amp config set default_project my-project
```

### Get Value

```bash
amp config get api_url
amp config get default_org
```

### List All

```bash
amp config list
```

### Available Keys

| Key | Description |
|-----|-------------|
| `api_url` | API server URL |
| `api_key_header` | Authentication header name |
| `api_key` | Authentication token |
| `default_org` | Default organization |
| `default_project` | Default project |

## Output Formats

### Table (default)

```bash
amp agents list
```

### JSON

```bash
amp agents list --output json
amp agents get my-agent --output json | jq '.status'
```

## Error Handling

The CLI provides clear error messages for common issues:

- **401 Unauthorized**: Invalid or expired token
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource doesn't exist
- **409 Conflict**: Resource already exists

## Configuration File

Config stored at `~/.amp/config.yaml`:

```yaml
api_url: https://api.example.com/api/v1
api_key_header: Authorization
api_key: Bearer eyJ...
default_org: my-org
default_project: my-project
```
