package github

type CopilotUsage struct {
	Day                   string          `json:"day"`
	TotalSuggestionsCount int             `json:"total_suggestions_count"`
	TotalAcceptancesCount int             `json:"total_acceptances_count"`
	TotalLinesSuggested   int             `json:"total_lines_suggested"`
	TotalLinesAccepted    int             `json:"total_lines_accepted"`
	TotalActiveUsers      int             `json:"total_active_users"`
	Breakdown             []LangBreakdown `json:"breakdown"`
}

type LangBreakdown struct {
	Language         string `json:"language"`
	Editor           string `json:"editor"`
	SuggestionsCount int    `json:"suggestions_count"`
	AcceptancesCount int    `json:"acceptances_count"`
	LinesSuggested   int    `json:"lines_suggested"`
	LinesAccepted    int    `json:"lines_accepted"`
	ActiveUsers      int    `json:"active_users"`
}

type SeatBreakdown struct {
	Total               int `json:"total"`
	AddedThisCycle      int `json:"added_this_cycle"`
	PendingInvitation   int `json:"pending_invitation"`
	PendingCancellation int `json:"pending_cancellation"`
	ActiveThisCycle     int `json:"active_this_cycle"`
	InactiveThisCycle   int `json:"inactive_this_cycle"`
}

type CopilotBilling struct {
	SeatBreakdown         SeatBreakdown `json:"seat_breakdown"`
	SeatManagementSetting string        `json:"seat_management_setting"`
	PublicCodeSuggestions string        `json:"public_code_suggestions"`
}

type CopilotBillingSeat struct {
	Assignee                CopilotAssignee     `json:"assignee"`
	AssigningTeam           CopilotAssigneeTeam `json:"assigning_team,omitempty"`
	CreatedAt               string              `json:"created_at,omitempty"`
	UpdatedAt               string              `json:"updated_at,omitempty"`
	LastActivityAt          string              `json:"last_activity_at,omitempty"`
	LastActivityEditor      string              `json:"last_activity_editor,omitempty"`
	PendingCancellationDate string              `json:"pending_cancellation_date,omitempty"`
}

type CopilotAssignee struct {
	AvatarURL         string `json:"avatar_url,omitempty"`
	EventsURL         string `json:"events_url,omitempty"`
	FollowersURL      string `json:"followers_url,omitempty"`
	FollowingURL      string `json:"following_url,omitempty"`
	GistsURL          string `json:"gists_url,omitempty"`
	GravatarID        string `json:"gravatar_id,omitempty"`
	HtmlURL           string `json:"html_url,omitempty"`
	Id                int    `json:"id,omitempty"`
	Login             string `json:"login,omitempty"`
	NodeId            string `json:"node_id,omitempty"`
	OrganizationsURL  string `json:"organizations_url,omitempty"`
	ReceivedEventsURL string `json:"received_events_url,omitempty"`
	ReposURL          string `json:"repos_url,omitempty"`
	SiteAdmin         bool   `json:"site_admin,omitempty"`
	StarredURL        string `json:"starred_url,omitempty"`
	SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
	Type              string `json:"type,omitempty"`
	URL               string `json:"url,omitempty"`
}

type CopilotAssigneeTeam struct {
	Id                   int    `json:"id,omitempty"`
	NodeId               string `json:"node_id,omitempty"`
	URL                  string `json:"url,omitempty"`
	HtmlURL              string `json:"html_url,omitempty"`
	Name                 string `json:"name,omitempty"`
	Slug                 string `json:"slug,omitempty"`
	Description          string `json:"description,omitempty"`
	Privacy              string `json:"privacy,omitempty"`
	NotificationSettings string `json:"notification_setting,omitempty"`
	Permission           string `json:"permission,omitempty"`
	Email                string `json:"email,omitempty"`
	MembersURL           string `json:"members_url,omitempty"`
	RepositoriesURL      string `json:"repositories_url,omitempty"`
	Parent               string `json:"parent,omitempty"`
}

type CopilotBillingSeats struct {
	TotalSeats int                  `json:"total_seats,omitempty"`
	Seats      []CopilotBillingSeat `json:"seats,omitempty"`
}
