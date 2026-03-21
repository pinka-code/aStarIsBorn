package model

type User struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"` // User or Organization
}

type Owner = User

type Contributor struct {
	User
	Contributions int `json:"contributions"`
}
