# Email Verifier

A Go tool that checks domain email security configurations (MX, SPF, DMARC records).

## ⚠️ Important Limitation

**Only verifies DOMAINS** (like `gmail.com`) - **NOT individual email addresses** (like `user@gmail.com`).

Individual email verification is impossible due to privacy/security restrictions.

## Quick Start

### API Server (for web integration)
```bash
go mod tidy
go build -o email-verifier-api.exe api_server.go
./email-verifier-api.exe
```
Access: `http://localhost:8080`

### CLI Tool (for command line)
```bash
cd cli
go build -o email-verifier-cli.exe main.go
echo "gmail.com" | ./email-verifier-cli.exe
```

## What It Checks

- **MX Records**: Can domain receive emails?
- **SPF Records**: Email spoofing protection
- **DMARC Records**: Email security policies

## API Endpoints

```
GET  /api/v1/health
GET  /api/v1/verify/domain/{domain}
POST /api/v1/verify/domains
```

## Example Usage

**API:**
```bash
curl http://localhost:8080/api/v1/verify/domain/gmail.com
```

**CLI:**
```bash
echo "gmail.com" | ./email-verifier-cli.exe
```

## Output Format

```json
{
  "domain": "gmail.com",
  "hasMX": true,
  "hasSPF": true,
  "spfRecord": "v=spf1 include:_spf.google.com ~all",
  "hasDMARC": true,
  "dmarcRecord": "v=DMARC1; p=reject"
}
```

## Use Cases

- Domain security auditing
- Email infrastructure validation
- Compliance checking
- Bulk domain analysis

**Note**: Cannot verify if individual users exist - this is technically impossible.