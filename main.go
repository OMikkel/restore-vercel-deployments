package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/omikkel/restore-vercel-deployments/internal/config"
	"github.com/omikkel/restore-vercel-deployments/internal/logger"
	"github.com/omikkel/restore-vercel-deployments/internal/vercel"
)

type Output struct {
	GeneratedAt        string                                          `yaml:"generated_at"`
	Teams              []vercel.VercelTeam                             `yaml:"teams"`
	ProjectsPerTeam    map[string][]vercel.VercelProject               `yaml:"projects_per_team"`
	DeletedDeployments map[string]map[string][]vercel.VercelDeployment `yaml:"deleted_deployments"`
}

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("[ERROR]", err)
		fmt.Println("[ERROR] Set VERCEL_API_TOKEN in your environment or create a .env file")
		os.Exit(1)
	}

	log := logger.NewLogger(cfg.LogLevel)
	log.Info("Starting Vercel Deployments Restorer")

	vercelAPI := vercel.NewVercelAPI(log, cfg.APIURL, cfg.APIToken)

	output := Output{
		GeneratedAt:        time.Now().Format(time.RFC3339),
		ProjectsPerTeam:    make(map[string][]vercel.VercelProject),
		DeletedDeployments: make(map[string]map[string][]vercel.VercelDeployment),
	}

	teams, _ := vercelAPI.GetTeams(0)
	output.Teams = teams
	log.Info("Found teams:", len(teams))

	for _, team := range teams {
		projects, _ := vercelAPI.GetProjects(team.ID, 0)
		output.ProjectsPerTeam[team.ID] = projects
		log.Info("Team:", team.Name, "("+team.ID+")", "Projects:", len(projects))

		for _, project := range projects {
			log.Info(" - Project:", project.Name, "ID:", project.ID)

			deletedDeployments, _ := vercelAPI.GetDeletedDeploymentsFromProject(team.ID, project.ID, 0)
			if output.DeletedDeployments[team.ID] == nil {
				output.DeletedDeployments[team.ID] = make(map[string][]vercel.VercelDeployment)
			}
			output.DeletedDeployments[team.ID][project.ID] = deletedDeployments
			log.Info("   Deleted Deployments:", len(deletedDeployments))

			for _, deployment := range deletedDeployments {
				log.Info("     - Deployment ID:", deployment.ID, "Branch:", deployment.Branch, "DeletedAt:", deployment.DeletedAt, "ByRetention:", deployment.DeletedByRetention)
				vercelAPI.RestoreDeploymentByID(team.ID, deployment.ID)

				log.Debug("Sleeping to avoid rate limits...")
				time.Sleep(cfg.RestoreCooldown)
			}
			log.Info("   Finished restoring " + fmt.Sprint(len(deletedDeployments)) + " deleted deployments for project " + project.Name)
		}
	}

	log.Info("Finished restoring deleted deployments for all projects and teams")
	log.Info("Writing processed deployments to file: ", ".out/deployment_overview.yaml")

	// write output to yaml file
	outputPath := filepath.Join(".out", "deployment_overview.yaml")
	os.MkdirAll(".out", os.ModePerm)
	log.Debug("Writing output to " + outputPath)

	formatted, _ := yaml.Marshal(output)
	os.WriteFile(outputPath, formatted, 0644)

	log.Info("Finished Vercel Deployments Restorer")
}
