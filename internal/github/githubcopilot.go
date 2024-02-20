package github

import (
	"encoding/json"
	"fmt"
)

func (c *GitHubClient) GetCopilotUsage() ([]CopilotUsage, error) {
	res, err := c.get("copilot/usage")
	if err != nil {
		return nil, fmt.Errorf("failed to get copilot usage data: %v", err)
	}

	var usage []CopilotUsage

	if err := json.NewDecoder(res.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return usage, nil
}

func (c *GitHubClient) GetBilling() (CopilotBilling, error) {
	res, err := c.get("copilot/billing")
	if err != nil {
		return CopilotBilling{}, fmt.Errorf("failed to get copilot billing data: %v", err)
	}

	var billing CopilotBilling

	if err := json.NewDecoder(res.Body).Decode(&billing); err != nil {
		return CopilotBilling{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return billing, nil
}

func (c *GitHubClient) GetBillingSeats() (CopilotBillingSeats, error) {
	//TODO: Response is paginated, so we need to handle that
	// see https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api?apiVersion=2022-11-28
	res, err := c.get("copilot/billing/seats")
	if err != nil {
		return CopilotBillingSeats{}, fmt.Errorf("failed to get copilot billing data: %v", err)
	}

	var billing CopilotBillingSeats

	if err := json.NewDecoder(res.Body).Decode(&billing); err != nil {
		return CopilotBillingSeats{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return billing, nil
}
