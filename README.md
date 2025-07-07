# copilot-metrics-exporter
This is simple GitHub copilot prometheus metrics exporter

### Description
GitHub Copilot is an AI pair programmer that helps you write code faster. GitHub Copilot draws context from the code you’re working on, 
suggesting whole lines or entire functions. It is a tool that helps developers write code faster and with fewer errors. 

This exporter is designed to collect and expose key metrics from GitHub Copilot to Prometheus. 
The metrics are collected from the GitHub Copilot API and are exposed in a format that can be scraped by Prometheus. 
The exporter is designed to be run as a standalone service and can be deployed as a containerized application.

### Key metrics
All metrics are provided by the last available day in the GitHub Copilot API. 

- Acceptance Rate: This metric represents the ratio of accepted lines to the total lines suggested by GitHub Copilot. This rate is an indicator of the relevance and usefulness of Copilot's suggestions.
- Total Suggestions This chart illustrates the total number of code suggestions made by GitHub Copilot. It offers a view of the tool's activity and its engagement with users over time.
- Total Acceptances: This visualization focuses on the total number of suggestions accepted by users.
- Total Lines Suggested: Showcases the total number of lines of code suggested by GitHub Copilot. This gives an idea of the volume of code generation and assistance provided.
- Total Lines Accepted: As the name says, the total lines of code accepted by users (full acceptances) offering insights into how much of the suggested code is actually being utilized incorporated to the codebase.
- Total Active Users: Represents the number of active users engaging with GitHub Copilot. This helps in understanding the user base growth and adoption rate.
- Languages breakdown

### How to run
- Create `config.yaml` file in the directory `.copilot-exporter` of the home directory or near the executable file
- File should contain the following content
```yaml
org: myorgname
pat: ghp_mygithubpat
```
- Default port is `9080` and can be changed by setting it in the config file
```yaml
# example
port: 9888
```
- If your Github organisation is an enterprise, set the `is_enterprise` flag to `yes`. By default, it is `no`
```yaml 
# example
is_enterprise: yes
```

or set the environment variables `GITHUB_ORG` `GITHUB_TOKEN` `GITHUB_IS_ENTERPRISE` and `PORT`
```bash
export GITHUB_ORG=myorgname
export GITHUB_TOKEN=ghp_mygithubpat
export PORT=9888
```

and run the executable file `copilot-metrics-exporter` or `go run ./cmd/copilot-exporter/main.go`

### Metrics
Available metrics are:
- github_copilot_total_acceptances_count 
- github_copilot_total_active_users 
- github_copilot_total_lines_accepted 
- github_copilot_total_lines_suggested 
- github_copilot_total_seats_occupied 
- github_copilot_total_suggestions_count
- github_copilot_suggestions_count_breakdown (labels: language, editor)
- github_copilot_lines_suggested_breakdown (labels: language, editor)
- github_copilot_lines_accepted_breakdown (labels: language, editor)
- github_copilot_acceptances_count_breakdown (labels: language, editor)
- github_copilot_total_chats
- github_copilot_total_chat_insertions
- github_copilot_total_chat_copies

### GitHub App Authentication
As an alternative to Personal Access Tokens, you can use GitHub Apps:

1. [Create a GitHub App](https://docs.github.com/en/developers/apps/creating-a-github-app) in your organization
2. Generate a private key for your app
3. Install the app to your organization and note the installation ID
4. Configure using one of the following methods:

#### Config file:
```yaml
org: myorgname
# Either use PAT
# pat: ghp_mygithubpat

# Or use GitHub App authentication
github_app:
  app_id: 12345  
  installation_id: 987654
  private_key_path: /path/to/private-key.pem
  # OR private_key directly (less secure)
  # private_key: |
  #   -----BEGIN RSA PRIVATE KEY-----
  #   ...
  #   -----END RSA PRIVATE KEY-----

# Optional TLS configuration
tls:
  enabled: true
  cert_file: /path/to/certificate.pem
  key_file: /path/to/private-key.pem
```

### Security Notes
- When TLS is enabled, the exporter will serve HTTPS on the configured port
- Ensure your certificate files are properly secured with appropriate file permissions
- The TLS implementation uses TLS 1.2 as the minimum version for security
- For production use, consider using certificates from a trusted Certificate Authority (such as "Let's Encrypt")
