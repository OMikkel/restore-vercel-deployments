package vercel

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/omikkel/restore-vercel-deployments/internal/logger"
	"github.com/omikkel/restore-vercel-deployments/internal/utils"
)

type VercelAPI struct {
	Logger     *logger.Logger
	API_URL    string
	API_TOKEN  string
	HTTPClient *http.Client
}

type VercelTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VercelProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VercelDeployment struct {
	ID                 string `json:"id"`
	Branch             string `json:"branch"`
	CommitSHA          string `json:"commit_sha"`
	DeletedAt          int    `json:"deleted_at"`
	DeletedByRetention bool   `json:"deleted_by_retention"`
}

type VercelPagination struct {
	Count int `json:"count"`
	Next  int `json:"next"`
	Prev  int `json:"prev"`
}

func NewVercelAPI(logger *logger.Logger, apiUrl string, apiToken string) *VercelAPI {
	return &VercelAPI{
		Logger:     logger,
		API_URL:    apiUrl,
		API_TOKEN:  apiToken,
		HTTPClient: &http.Client{},
	}
}

func (v *VercelAPI) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+v.API_TOKEN)
	return v.HTTPClient.Do(req)
}

func (v *VercelAPI) requestWithPagination(req *http.Request) (map[string]interface{}, *http.Response, VercelPagination, error) {
	res, err := v.doRequest(req)
	if err != nil {
		return nil, nil, VercelPagination{}, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	pagination := VercelPagination{}
	if pagMap, ok := result["pagination"].(map[string]interface{}); ok {
		pagination.Count = int(pagMap["count"].(float64))
		if pagMap["next"] != nil {
			pagination.Next = int(pagMap["next"].(float64))
		}
		if pagMap["prev"] != nil {
			pagination.Prev = int(pagMap["prev"].(float64))
		}
	}

	return result, res, pagination, nil
}

func (v *VercelAPI) GetTeams(until int) ([]VercelTeam, error) {
	endpoint := utils.URLWithQueryParams(v.API_URL+"/v2/teams", map[string]string{
		"until": func() string {
			if until != 0 {
				return strconv.Itoa(until)
			}
			return ""
		}(),
	})

	// request vercel api
	if until != 0 {
		v.Logger.Debug("Getting teams until: " + strconv.Itoa(until))
	} else {
		v.Logger.Debug("Getting teams")
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	result, res, pagination, err := v.requestWithPagination(req)
	if err != nil {
		return nil, err
	}

	v.Logger.Debug("Response Status: " + res.Status)
	v.Logger.Debug("Response Body:")

	jsonBody, _ := json.MarshalIndent(result, "", "  ")
	v.Logger.Debug(string(jsonBody))

	teams := []VercelTeam{}
	if items, ok := result["teams"].([]interface{}); ok {
		for _, item := range items {
			if teamMap, ok := item.(map[string]interface{}); ok {
				team := VercelTeam{
					ID:   teamMap["id"].(string),
					Name: teamMap["name"].(string),
				}
				teams = append(teams, team)
			}
		}
	}

	if pagination.Next != 0 {
		moreTeams, err := v.GetTeams(pagination.Next)
		if err != nil {
			return nil, err
		}
		teams = append(teams, moreTeams...)
	}

	return teams, nil
}

func (v *VercelAPI) GetProjects(teamID string, until int) ([]VercelProject, error) {
	endpoint := utils.URLWithQueryParams(v.API_URL+"/v10/projects", map[string]string{
		"teamId": teamID,
		"until": func() string {
			if until != 0 {
				return strconv.Itoa(until)
			}
			return ""
		}(),
	})

	if until != 0 {
		v.Logger.Debug("Getting projects for team ID: " + teamID + " until: " + strconv.Itoa(until))
	} else {
		v.Logger.Debug("Getting projects for team ID: " + teamID)
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	result, res, pagination, err := v.requestWithPagination(req)
	if err != nil {
		return nil, err
	}

	v.Logger.Debug("Response Status: " + res.Status)
	v.Logger.Debug("Response Body:")

	jsonBody, _ := json.MarshalIndent(result, "", "  ")
	v.Logger.Debug(string(jsonBody))

	projects := []VercelProject{}
	if items, ok := result["projects"].([]interface{}); ok {
		for _, item := range items {
			if projectMap, ok := item.(map[string]interface{}); ok {
				project := VercelProject{
					ID:   projectMap["id"].(string),
					Name: projectMap["name"].(string),
				}
				projects = append(projects, project)
			}
		}
	}

	if pagination.Next != 0 {
		moreProjects, err := v.GetProjects(teamID, pagination.Next)
		if err != nil {
			return nil, err
		}
		projects = append(projects, moreProjects...)
	}

	return projects, nil
}

func (v *VercelAPI) GetDeletedDeploymentsFromProject(teamID string, projectID string, until int) ([]VercelDeployment, error) {
	endpoint := utils.URLWithQueryParams(v.API_URL+"/v6/deployments", map[string]string{
		"limit":     "100",
		"projectId": projectID,
		"teamId":    teamID,
		"state":     "DELETED",
		"until": func() string {
			if until != 0 {
				return strconv.Itoa(until)
			}
			return ""
		}(),
	})

	if until != 0 {
		v.Logger.Debug("Getting deleted deployments for project ID: " + projectID + " in team ID: " + teamID + " until: " + strconv.Itoa(until))
	} else {
		v.Logger.Debug("Getting deleted deployments for project ID: " + projectID + " in team ID: " + teamID)
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	result, res, pagination, err := v.requestWithPagination(req)
	if err != nil {
		return nil, err
	}

	v.Logger.Debug("Response Status: " + res.Status)
	v.Logger.Debug("Response Body:")

	jsonBody, _ := json.MarshalIndent(result, "", "  ")
	v.Logger.Debug(string(jsonBody))

	deletedDeployments := []VercelDeployment{}
	if items, ok := result["deployments"].([]interface{}); ok {
		for _, item := range items {
			if deploymentMap, ok := item.(map[string]interface{}); ok {
				deployment := VercelDeployment{
					ID:                 deploymentMap["uid"].(string),
					DeletedAt:          int(deploymentMap["deleted"].(float64)),
					DeletedByRetention: deploymentMap["softDeletedByRetention"].(bool),
				}
				if deploymentMap["meta"] != nil {
					if metaMap, ok := deploymentMap["meta"].(map[string]interface{}); ok {
						if metaMap["githubCommitRef"] != nil {
							deployment.Branch = metaMap["githubCommitRef"].(string)
						}
						if metaMap["githubCommitSha"] != nil {
							deployment.CommitSHA = metaMap["githubCommitSha"].(string)
						}
					}
				}
				deletedDeployments = append(deletedDeployments, deployment)
			}
		}
	}

	if pagination.Next != 0 {
		moreDeletedDeployments, err := v.GetDeletedDeploymentsFromProject(teamID, projectID, pagination.Next)
		if err != nil {
			return nil, err
		}
		deletedDeployments = append(deletedDeployments, moreDeletedDeployments...)
	}

	return deletedDeployments, nil
}

func (v *VercelAPI) RestoreDeploymentByID(teamID string, deploymentID string) error {
	endpoint := utils.URLWithQueryParams(v.API_URL+"/v1/projects/undelete-deployment/"+deploymentID, map[string]string{
		"teamId": teamID,
	})

	v.Logger.Debug("Restoring deployment ID: " + deploymentID + " in team ID: " + teamID)
	req, err := http.NewRequest("PATCH", endpoint, nil)
	req.Header.Add("Content-Type", "application/json")
	res, err := v.doRequest(req)
	if err != nil {
		return err
	}

	v.Logger.Debug("Got code: " + res.Status + ", Restored deployment ID: " + deploymentID)

	return nil
}
