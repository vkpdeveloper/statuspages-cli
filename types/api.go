package types

type Page struct {
	Id              string `json:"id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Name            string `json:"name"`
	PageDescription string `json:"page_description"`
	Headline        string `json:"headline"`
	Subdomain       string `json:"subdomain"`
	Domain          string `json:"domain"`
	Url             string `json:"url"`
}

type Component struct {
	Id          string `json:"id"`
	PageId      string `json:"page_id"`
	UpdatedAt   string `json:"updated_at"`
	Group       bool   `json:"group"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	Status      string `json:"status"`
}

type UserPermission struct {
	Data []struct {
		UserId string           `json:"user_id"`
		Pages  []PermissionPage `json:"pages"`
	} `json:"data"`
}

type PermissionPage struct {
	PageId             string `json:"page_id"`
	PageConfiguration  bool   `json:"page_configuration"`
	IncidentManager    bool   `json:"incident_manager"`
	MaintenanceManager bool   `json:"maintenance_manager"`
}
