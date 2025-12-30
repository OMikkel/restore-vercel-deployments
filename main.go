package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/omikkel/restore-vercel-deployments/internal/logger"
	"github.com/omikkel/restore-vercel-deployments/internal/vercel"
)

const (
	LOG_LEVEL        = logger.LevelInfo          // Change to LevelDebug for more verbose output
	VERCEL_API_URL   = "https://vercel.com/api"  // or your custom Vercel API URL
	VERCEL_API_TOKEN = "<YOUR_VERCEL_API_TOKEN>" // Set your Vercel API token here
	RESTORE_COOLDOWN = 250 * time.Millisecond    // Time to wait between restore requests to avoid rate limits
)

type Output struct {
	GeneratedAt        string                                          `yaml:"generated_at"`
	Teams              []vercel.VercelTeam                             `yaml:"teams"`
	ProjectsPerTeam    map[string][]vercel.VercelProject               `yaml:"projects_per_team"`
	DeletedDeployments map[string]map[string][]vercel.VercelDeployment `yaml:"deleted_deployments"`
}

func main() {
	logger := logger.NewLogger(LOG_LEVEL)
	logger.Info("Starting Vercel Deployments Restorer")

	vercelAPI := vercel.NewVercelAPI(logger, VERCEL_API_URL, VERCEL_API_TOKEN)

	output := Output{
		GeneratedAt:        time.Now().Format(time.RFC3339),
		ProjectsPerTeam:    make(map[string][]vercel.VercelProject),
		DeletedDeployments: make(map[string]map[string][]vercel.VercelDeployment),
	}

	teams, _ := vercelAPI.GetTeams(0)
	output.Teams = teams
	logger.Info("Found teams:", len(teams))

	for _, team := range teams {
		projects, _ := vercelAPI.GetProjects(team.ID, 0)
		output.ProjectsPerTeam[team.ID] = projects
		logger.Info("Team:", team.Name, "("+team.ID+")", "Projects:", len(projects))

		for _, project := range projects {
			logger.Info(" - Project:", project.Name, "ID:", project.ID)

			deletedDeployments, _ := vercelAPI.GetDeletedDeploymentsFromProject(team.ID, project.ID, 0)
			if output.DeletedDeployments[team.ID] == nil {
				output.DeletedDeployments[team.ID] = make(map[string][]vercel.VercelDeployment)
			}
			output.DeletedDeployments[team.ID][project.ID] = deletedDeployments
			logger.Info("   Deleted Deployments:", len(deletedDeployments))

			for _, deployment := range deletedDeployments {
				logger.Info("     - Deployment ID:", deployment.ID, "Branch:", deployment.Branch, "DeletedAt:", deployment.DeletedAt, "ByRetention:", deployment.DeletedByRetention)
				vercelAPI.RestoreDeploymentByID(team.ID, deployment.ID)

				logger.Debug("Sleeping for 250ms to avoid rate limits...")
				time.Sleep(RESTORE_COOLDOWN)
			}
			logger.Info("   Finished restoring " + fmt.Sprint(len(deletedDeployments)) + " deleted deployments for project " + project.Name)
		}
	}

	logger.Info("Finished restoring deleted deployments for all projects and teams")
	logger.Info("Writing processed deployments to file: ", ".out/deployment_overview.yaml")

	// write output to yaml file
	filepath := filepath.Join(".out", "deployment_overview.yaml")
	os.MkdirAll(".out", os.ModePerm)
	logger.Debug("Writing output to " + filepath)

	formatted, _ := yaml.Marshal(output)
	os.WriteFile(filepath, formatted, 0644)

	logger.Info("Finished Vercel Deployments Restorer")
}
