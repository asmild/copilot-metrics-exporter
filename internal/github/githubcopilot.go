package github

import (
	"encoding/json"
	"fmt"
)

// GetCopilotMetrics https://docs.github.com/en/enterprise-cloud@latest/rest/copilot/copilot-metrics?apiVersion=2022-11-28
func (c *Client) GetCopilotMetrics(since *string) ([]CopilotMetrics, error) {
	endpoint := fmt.Sprintf("copilot/metrics?since=%s", *since)

	res, err := c.get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get copilot usage data: %v", err)
	}

	var usage []CopilotMetrics

	if err := json.NewDecoder(res.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}
	return usage, nil
}

// GetBilling https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin/billing?apiVersion=2022-11-28#get-billing-usage-report-for-an-enterprise
func (c *Client) GetBilling() (CopilotBilling, error) {
	res, err := c.get("copilot/billing")
	//  https://api.github.com/enterprises/ENTERPRISE/settings/billing/usage
	if err != nil {
		return CopilotBilling{}, fmt.Errorf("failed to get copilot billing data: %v", err)
	}

	var billing CopilotBilling

	if err := json.NewDecoder(res.Body).Decode(&billing); err != nil {
		return CopilotBilling{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return billing, nil
}

// GetBillingSeats https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api?apiVersion=2022-11-28
func (c *Client) GetBillingSeats() (CopilotBillingSeats, error) {
	//TODO: Response is paginated, so we need to handle that
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
