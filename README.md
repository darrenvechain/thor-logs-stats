# thor-logs-stats

1. Get the logs (specify your AWS profile):
```html
aws logs filter-log-events --profile <aws-profile> \
    --log-group-name <log-group-name> \
    --start-time $(($(date +%s) - 3600))000 \
    --filter-pattern 'API Request' \
    --output json > logs_output.json
```

2. Run the script:
```bash
go run .
```
