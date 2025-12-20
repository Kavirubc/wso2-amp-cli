package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
)

// Executor handles command execution in interactive mode
type Executor struct {
	client *api.Client
}

// NewExecutor creates a new command executor
func NewExecutor() *Executor {
	client := api.NewClient(
		config.GetAPIURL(),
		config.GetAPIKeyHeader(),
		config.GetAPIKeyValue(),
	)
	return &Executor{client: client}
}

// Execute parses and runs a command, returning the output
func (e *Executor) Execute(input string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return ""
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	switch cmd {
	case "help", "?":
		return e.help()
	case "orgs":
		return e.orgs(args)
	case "projects":
		return e.projects(args)
	case "agents":
		return e.agents(args)
	case "config":
		return e.config(args)
	default:
		return ui.RenderError(fmt.Sprintf("Unknown command: %s. Type 'help' for available commands.", cmd))
	}
}

func (e *Executor) help() string {
	return `Available Commands:
  orgs list              List all organizations
  projects list          List projects
  projects get <name>    Get project details
  agents list            List agents
  config list            Show configuration
  config set <k> <v>     Set config value
  clear                  Clear screen
  exit                   Exit interactive mode`
}

func (e *Executor) orgs(args []string) string {
	if len(args) == 0 || args[0] != "list" {
		return ui.RenderWarning("Usage: orgs list")
	}

	orgs, err := e.client.ListOrganizations()
	if err != nil {
		return ui.RenderError(err.Error())
	}

	if len(orgs) == 0 {
		return ui.RenderWarning("No organizations found.")
	}

	headers := []string{"NAME", "CREATED AT"}
	rows := make([][]string, len(orgs))
	for i, org := range orgs {
		rows[i] = []string{org.Name, org.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	return ui.RenderTableWithTitle("Organizations", headers, rows)
}

func (e *Executor) projects(args []string) string {
	if len(args) == 0 {
		return ui.RenderWarning("Usage: projects list | projects get <name>")
	}

	org := config.GetDefaultOrg()
	if org == "" {
		return ui.RenderError("No default org set. Run: config set default_org <name>")
	}

	switch args[0] {
	case "list":
		projects, err := e.client.ListProjects(org)
		if err != nil {
			return ui.RenderError(err.Error())
		}
		if len(projects) == 0 {
			return ui.RenderWarning("No projects found.")
		}
		headers := []string{"NAME", "DISPLAY NAME", "CREATED AT"}
		rows := make([][]string, len(projects))
		for i, p := range projects {
			rows[i] = []string{p.Name, p.DisplayName, p.CreatedAt.Format("2006-01-02 15:04:05")}
		}
		return ui.RenderTableWithTitle(fmt.Sprintf("Projects in %s", org), headers, rows)

	case "get":
		if len(args) < 2 {
			return ui.RenderWarning("Usage: projects get <name>")
		}
		project, err := e.client.GetProject(org, args[1])
		if err != nil {
			return ui.RenderError(err.Error())
		}
		var buf bytes.Buffer
		buf.WriteString(ui.TitleStyle.Render(fmt.Sprintf("Project: %s", project.Name)))
		buf.WriteString("\n\n")
		buf.WriteString(fmt.Sprintf("  %s  %s\n", ui.KeyStyle.Render("Name:"), project.Name))
		buf.WriteString(fmt.Sprintf("  %s  %s\n", ui.KeyStyle.Render("Display Name:"), project.DisplayName))
		buf.WriteString(fmt.Sprintf("  %s  %s\n", ui.KeyStyle.Render("Organization:"), project.OrgName))
		buf.WriteString(fmt.Sprintf("  %s  %s\n", ui.KeyStyle.Render("Created At:"), project.CreatedAt.Format("2006-01-02 15:04:05")))
		return buf.String()

	default:
		return ui.RenderWarning("Usage: projects list | projects get <name>")
	}
}

func (e *Executor) agents(args []string) string {
	if len(args) == 0 || args[0] != "list" {
		return ui.RenderWarning("Usage: agents list")
	}

	org := config.GetDefaultOrg()
	project := config.GetDefaultProject()
	if org == "" || project == "" {
		return ui.RenderError("Set defaults first: config set default_org/default_project <name>")
	}

	agents, err := e.client.ListAgents(org, project)
	if err != nil {
		return ui.RenderError(err.Error())
	}

	if len(agents) == 0 {
		return ui.RenderWarning("No agents found.")
	}

	headers := []string{"NAME", "DISPLAY NAME", "STATUS"}
	rows := make([][]string, len(agents))
	for i, a := range agents {
		rows[i] = []string{a.Name, a.DisplayName, ui.StatusCell(a.Status)}
	}
	return ui.RenderTableWithTitle(fmt.Sprintf("Agents in %s/%s", org, project), headers, rows)
}

func (e *Executor) config(args []string) string {
	if len(args) == 0 {
		return ui.RenderWarning("Usage: config list | config set <key> <value>")
	}

	switch args[0] {
	case "list":
		var buf bytes.Buffer
		buf.WriteString(ui.TitleStyle.Render("Configuration"))
		buf.WriteString("\n\n")
		buf.WriteString(fmt.Sprintf("  %s  %s\n", ui.KeyStyle.Render("api_url:"), config.GetAPIURL()))
		buf.WriteString(fmt.Sprintf("  %s  %s\n", ui.KeyStyle.Render("default_org:"), valueOr(config.GetDefaultOrg(), "(not set)")))
		buf.WriteString(fmt.Sprintf("  %s  %s\n", ui.KeyStyle.Render("default_project:"), valueOr(config.GetDefaultProject(), "(not set)")))
		return buf.String()

	case "set":
		if len(args) < 3 {
			return ui.RenderWarning("Usage: config set <key> <value>")
		}
		if err := config.Set(args[1], args[2]); err != nil {
			return ui.RenderError(err.Error())
		}
		return ui.RenderSuccess(fmt.Sprintf("Set %s = %s", args[1], args[2]))

	case "get":
		if len(args) < 2 {
			return ui.RenderWarning("Usage: config get <key>")
		}
		val := config.Get(args[1])
		if val == "" {
			return fmt.Sprintf("%s: (not set)", args[1])
		}
		encoder := json.NewEncoder(&bytes.Buffer{})
		encoder.SetIndent("", "")
		return fmt.Sprintf("%s: %s", args[1], val)

	default:
		return ui.RenderWarning("Usage: config list | config set <key> <value>")
	}
}

func valueOr(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}
