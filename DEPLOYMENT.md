# 🚀 Email Verifier API - Deployment Guide

## 📋 Overview

This API allows you to verify domain email configurations from your frontend applications. It provides REST endpoints to check MX, SPF, and DMARC records for domains.

## ⚠️ Important Note

**This API only verifies DOMAINS** (like `gmail.com`) **NOT individual email addresses** (like `user@gmail.com`). Individual email verification is impossible due to privacy and security restrictions.

## 🛠️ Local Development

### Prerequisites
- Go 1.19 or higher
- Git

### Steps

1. **Install dependencies:**
```bash
go mod tidy
```

2. **Run the API server:**
```bash
go run api_server.go
```

3. **Open the frontend example:**
   - Open `frontend_example.html` in your browser
   - The API will be running on `http://localhost:8080`

## 🐳 Docker Deployment

### Build and Run with Docker

```bash
# Build the Docker image
docker build -t email-verifier-api .

# Run the container
docker run -p 8080:8080 email-verifier-api
```

### Using Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'
services:
  email-verifier-api:
    build: .
    ports:
      - "8080:8080"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

Then run:
```bash
docker-compose up -d
```

## ☁️ Cloud Deployment Options

### 1. Heroku Deployment

Create `Procfile`:
```
web: ./email-verifier-api
```

Deploy:
```bash
# Login to Heroku
heroku login

# Create app
heroku create your-email-verifier-api

# Set buildpack
heroku buildpacks:set heroku/go

# Deploy
git add .
git commit -m "Deploy email verifier API"
git push heroku main
```

### 2. Railway Deployment

1. Connect your GitHub repo to Railway
2. Railway will auto-detect the Go application
3. Set environment variables if needed
4. Deploy automatically

### 3. Vercel Deployment

Create `vercel.json`:
```json
{
  "builds": [
    {
      "src": "api_server.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "/api_server.go"
    }
  ]
}
```

### 4. DigitalOcean App Platform

Create `.do/app.yaml`:
```yaml
name: email-verifier-api
services:
- name: api
  source_dir: /
  github:
    repo: your-username/your-repo
    branch: main
  run_command: ./email-verifier-api
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  http_port: 8080
  health_check:
    http_path: /api/v1/health
```

## 📡 API Endpoints

### Health Check
```
GET /api/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "time": "2024-10-12T10:30:00Z"
}
```

### Verify Single Domain
```
GET /api/v1/verify/domain/{domain}
```

**Example:**
```bash
curl http://localhost:8080/api/v1/verify/domain/gmail.com
```

**Response:**
```json
{
  "domain": "gmail.com",
  "hasMX": true,
  "hasSPF": true,
  "spfRecord": "v=spf1 include:_spf.google.com ~all",
  "hasDMARC": true,
  "dmarcRecord": "v=DMARC1; p=reject; ruf=mailto:mailauth-reports@google.com",
  "checkedAt": "2024-10-12T10:30:00Z"
}
```

### Verify Multiple Domains
```
POST /api/v1/verify/domains
Content-Type: application/json

{
  "domains": ["gmail.com", "github.com", "badexample.com"]
}
```

**Response:**
```json
{
  "results": [
    {
      "domain": "gmail.com",
      "hasMX": true,
      "hasSPF": true,
      "spfRecord": "v=spf1 include:_spf.google.com ~all",
      "hasDMARC": true,
      "dmarcRecord": "v=DMARC1; p=reject; ruf=mailto:mailauth-reports@google.com",
      "checkedAt": "2024-10-12T10:30:00Z"
    }
  ],
  "total": 1
}
```

## 🌐 Frontend Integration

### JavaScript Example

```javascript
const API_BASE_URL = 'https://your-api-domain.com/api/v1';

async function verifyDomain(domain) {
    try {
        const response = await fetch(`${API_BASE_URL}/verify/domain/${domain}`);
        const data = await response.json();
        
        if (response.ok) {
            console.log('Domain verification result:', data);
            return data;
        } else {
            console.error('Verification failed:', data.message);
        }
    } catch (error) {
        console.error('API Error:', error);
    }
}

// Usage
verifyDomain('gmail.com').then(result => {
    if (result.hasMX) {
        console.log('✅ Domain can receive emails');
    } else {
        console.log('❌ Domain cannot receive emails');
    }
});
```

### React Example

```jsx
import { useState } from 'react';

function DomainVerifier() {
    const [domain, setDomain] = useState('');
    const [result, setResult] = useState(null);
    const [loading, setLoading] = useState(false);

    const verifyDomain = async () => {
        setLoading(true);
        try {
            const response = await fetch(`/api/v1/verify/domain/${domain}`);
            const data = await response.json();
            setResult(data);
        } catch (error) {
            console.error('Error:', error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <input 
                value={domain} 
                onChange={(e) => setDomain(e.target.value)}
                placeholder="Enter domain (e.g., gmail.com)"
            />
            <button onClick={verifyDomain} disabled={loading}>
                {loading ? 'Verifying...' : 'Verify Domain'}
            </button>
            
            {result && (
                <div>
                    <h3>Results for {result.domain}:</h3>
                    <p>MX Records: {result.hasMX ? '✅' : '❌'}</p>
                    <p>SPF Records: {result.hasSPF ? '✅' : '❌'}</p>
                    <p>DMARC Records: {result.hasDMARC ? '✅' : '❌'}</p>
                </div>
            )}
        </div>
    );
}
```

## 🔒 Security Considerations

### Rate Limiting
Consider adding rate limiting to prevent abuse:

```go
import "golang.org/x/time/rate"

var limiter = rate.NewLimiter(10, 50) // 10 requests per second, burst of 50

func rateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### CORS Configuration
For production, restrict CORS to your frontend domain:

```go
c := cors.New(cors.Options{
    AllowedOrigins: []string{"https://yourdomain.com"},
    AllowedMethods: []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders: []string{"Content-Type"},
})
```

### Environment Variables
Use environment variables for configuration:

```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
```

## 📊 Monitoring

### Health Check Endpoint
The API includes a health check endpoint that you can use with monitoring services:

```bash
curl http://localhost:8080/api/v1/health
```

### Logging
Add structured logging for better monitoring:

```go
import "github.com/sirupsen/logrus"

logrus.WithFields(logrus.Fields{
    "domain": domain,
    "hasMX":  result.HasMX,
}).Info("Domain verification completed")
```

## 🚨 Limitations

1. **Cannot verify individual email addresses** - This is a technical limitation, not a bug
2. **DNS resolution depends on network** - Some domains may be unreachable
3. **Rate limiting by DNS providers** - Bulk operations may hit rate limits
4. **Timeout constraints** - DNS lookups have timeouts

## 🆘 Troubleshooting

### Common Issues

1. **CORS errors in browser:**
   - Make sure CORS is configured for your frontend domain
   - Check that the API server is running

2. **DNS resolution failures:**
   - Check internet connectivity
   - Some corporate networks block DNS queries

3. **Port 8080 already in use:**
   - Change the port in the code or kill the existing process

### Debug Mode
Run with debug logging:

```bash
export LOG_LEVEL=debug
go run api_server.go
```

## 📝 License

[Add your license here]

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request