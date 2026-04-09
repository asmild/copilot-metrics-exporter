package github

// UsageReport represents the new Copilot usage metrics report format
// returned by the /copilot/metrics/reports/* endpoints (apiVersion 2026-03-10).
type UsageReport struct {
	Day                          string                 `json:"day"`
	EnterpriseID                 string                 `json:"enterprise_id,omitempty"`
	DailyActiveUsers             int                    `json:"daily_active_users"`
	DailyActiveCLIUsers          int                    `json:"daily_active_cli_users"`
	WeeklyActiveUsers            int                    `json:"weekly_active_users"`
	MonthlyActiveUsers           int                    `json:"monthly_active_users"`
	MonthlyActiveChatUsers       int                    `json:"monthly_active_chat_users"`
	MonthlyActiveAgentUsers      int                    `json:"monthly_active_agent_users"`
	UserInitiatedInteractionCount int                   `json:"user_initiated_interaction_count"`
	CodeGenerationActivityCount  int                    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount  int                    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum         int                    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum      int                    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                  int                    `json:"loc_added_sum"`
	LocDeletedSum                int                    `json:"loc_deleted_sum"`
	TotalsByIDE                  []IDETotals            `json:"totals_by_ide"`
	TotalsByFeature              []FeatureTotals        `json:"totals_by_feature"`
	TotalsByLanguageFeature      []LanguageFeatureTotals `json:"totals_by_language_feature"`
	PullRequests                 PullRequestTotals      `json:"pull_requests"`
}

// IDETotals represents per-IDE usage totals.
type IDETotals struct {
	IDE                           string `json:"ide"`
	UserInitiatedInteractionCount int    `json:"user_initiated_interaction_count"`
	CodeGenerationActivityCount   int    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount   int    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum          int    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum       int    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                   int    `json:"loc_added_sum"`
	LocDeletedSum                 int    `json:"loc_deleted_sum"`
}

// FeatureTotals represents per-feature usage totals (e.g. code_completion, chat_panel_agent_mode).
type FeatureTotals struct {
	Feature                       string `json:"feature"`
	UserInitiatedInteractionCount int    `json:"user_initiated_interaction_count"`
	CodeGenerationActivityCount   int    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount   int    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum          int    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum       int    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                   int    `json:"loc_added_sum"`
	LocDeletedSum                 int    `json:"loc_deleted_sum"`
}

// LanguageFeatureTotals represents per-language-per-feature usage totals.
type LanguageFeatureTotals struct {
	Language                    string `json:"language"`
	Feature                     string `json:"feature"`
	CodeGenerationActivityCount int    `json:"code_generation_activity_count"`
	CodeAcceptanceActivityCount int    `json:"code_acceptance_activity_count"`
	LocSuggestedToAddSum        int    `json:"loc_suggested_to_add_sum"`
	LocSuggestedToDeleteSum     int    `json:"loc_suggested_to_delete_sum"`
	LocAddedSum                 int    `json:"loc_added_sum"`
	LocDeletedSum               int    `json:"loc_deleted_sum"`
}

// PullRequestTotals represents pull request usage totals.
type PullRequestTotals struct {
	TotalReviewed              int `json:"total_reviewed"`
	TotalCreated               int `json:"total_created"`
	TotalCreatedByCopilot      int `json:"total_created_by_copilot"`
	TotalReviewedByCopilot     int `json:"total_reviewed_by_copilot"`
	TotalMerged                int `json:"total_merged"`
	TotalSuggestions           int `json:"total_suggestions"`
	TotalAppliedSuggestions    int `json:"total_applied_suggestions"`
}