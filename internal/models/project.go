package models

type Project struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TechStack   []string `json:"techStack"`
	GitHubUrl   string   `json:"gitHubUrl"`
	LiveUrl     string   `json:"liveUrl"`
	Image       string   `json:"image"`
	CreatedAt   string   `json:"createdAt"`
}
