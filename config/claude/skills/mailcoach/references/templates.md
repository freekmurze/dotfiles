# Templates

Templates are reusable HTML layouts that can be applied to campaigns.

## List templates

```bash
mailcoach list-templates
mailcoach list-templates --filter-search="newsletter" --json
```

## Create a template

```bash
mailcoach create-template --input '{
  "name": "Company Newsletter",
  "html": "<html><body><h1>{{title}}</h1><div>{{content}}</div></body></html>"
}'
```

## Show, update, delete

```bash
mailcoach show-template --template=<uuid>

mailcoach update-template --template=<uuid> --input '{
  "name": "Updated Template",
  "html": "<html><body><h1>New Layout</h1>{{content}}</body></html>"
}'

mailcoach delete-template --template=<uuid>
```

## Using templates with campaigns

When creating a campaign, pass the template UUID to use its HTML as the base:

```bash
mailcoach create-campaign --input '{
  "name": "April Newsletter",
  "email_list_uuid": "<list-uuid>",
  "template_uuid": "<template-uuid>"
}'
```

The campaign's HTML content inherits from the template and can be customized via `update-campaign`.
