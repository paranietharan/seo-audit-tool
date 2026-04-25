# SEO Audit tool in go

```bash
go run cmd/server/main.go
```

## Analyze the website

```bash
postman request POST 'localhost:8080/audit' \
  --header 'Content-Type: application/json' \
  --body '{
    "url":"https://paranietharan.vercel.app"
}'
```

## Get the audit report

```bash
postman request 'localhost:8080/report/a4eade58-bf9b-4a27-9a6e-46a3cc2f3486'
```