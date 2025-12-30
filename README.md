# Restore Vercel Deployments

A CLI tool to restore deleted deployments across all your Vercel teams and projects.

## Why?

Vercel's retention policy automatically deletes old deployments. This tool helps you restore those deleted deployments in bulk, which is useful for compliance, auditing, or recovering rollback options.

## Features

- 🔄 Automatically restores all deleted deployments across teams and projects
- 📊 Generates a YAML report of all teams, projects, and restored deployments
- ⏱️ Built-in rate limiting to avoid API throttling
- 📝 Configurable logging levels

## Prerequisites

- Go 1.23 or later
- A Vercel API token with appropriate permissions

## Installation

```bash
git clone https://github.com/omikkel/restore-vercel-deployments.git
cd restore-vercel-deployments
go mod download
```

## Configuration

Edit `main.go` and set your configuration:

```go
const (
    LOG_LEVEL        = logger.LevelInfo          // Options: LevelDebug, LevelInfo, LevelError, LevelDisabled
    VERCEL_API_URL   = "https://vercel.com/api"  // Default Vercel API URL
    VERCEL_API_TOKEN = "<YOUR_VERCEL_API_TOKEN>" // Your Vercel API token
    RESTORE_COOLDOWN = 250 * time.Millisecond    // Time to wait between restore requests to avoid rate limits
)
```

### Getting a Vercel API Token

1. Go to your [Vercel Account Settings](https://vercel.com/account/tokens)
2. Click "Create Token"
3. Give it a descriptive name and select the appropriate scope
4. Copy the token and paste it in `main.go`

## Usage

```bash
make run-dev
or
go run main.go
```

Or build and run:

```bash
make run
or
go build -o .out/restore-vercel-deployments
./.out/restore-vercel-deployments
```

## Output

The tool generates a YAML file at `.out/deployment_overview.yaml` containing:

- Timestamp of when the tool was run
- List of all teams
- Projects per team
- Deleted deployments that were restored (with deployment ID, branch, commit SHA, and deletion info)

Example output structure:

```yaml
generated_at: "2025-01-01T12:00:00Z"
teams:
  - id: team_xxxxx
    name: My Team
projects_per_team:
  team_xxxxx:
    - id: prj_xxxxx
      name: my-project
deleted_deployments:
  team_xxxxx:
    prj_xxxxx:
      - id: dpl_xxxxx
        branch: main
        commit_sha: abc123
        deleted_at: 1704067200
        deleted_by_retention: true
```

## Log Levels

| Level | Description |
|-------|-------------|
| `LevelDebug` | Verbose output including API responses |
| `LevelInfo` | Standard progress information |
| `LevelError` | Only error messages |
| `LevelDisabled` | No output |

## License

MIT

## Use of AI

This README was generated with the assistance of AI tools. While efforts have been made to ensure accuracy, please review the content for any potential errors or omissions.

The general script and logic were created by the author, with AI used to enhance documentation and suggest improvements.
