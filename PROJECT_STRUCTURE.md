# 📂 Project Structure

Your Email Verifier project now has two separate versions:

## 🚀 **API Server** (Recommended for Frontend Integration)
- **File**: `api_server.go`
- **Purpose**: REST API for web applications
- **Build**: `go build -o email-verifier-api.exe api_server.go`
- **Run**: `./email-verifier-api.exe`
- **URL**: `http://localhost:8080`

## 💻 **CLI Tool** (For Command Line Usage)
- **File**: `cli/main.go`
- **Purpose**: Command-line interface for batch processing
- **Build**: `cd cli && go build -o email-verifier-cli.exe main.go`
- **Usage**: `echo "gmail.com" | email-verifier-cli.exe`

## ⚡ **Quick Start**

### Option 1: Use the Deploy Script
```bash
# Run the interactive build script
deploy.bat

# Choose:
# 1 = API Server only
# 2 = CLI Tool only  
# 3 = Both
```

### Option 2: Manual Build

**For API Server:**
```bash
go mod tidy
go build -o email-verifier-api.exe api_server.go
./email-verifier-api.exe
```

**For CLI Tool:**
```bash
cd cli
go build -o email-verifier-cli.exe main.go
echo "gmail.com" | ./email-verifier-cli.exe
```

## 🌐 **Frontend Integration**

Open `frontend_example.html` in your browser to see the API in action, or use the React component `EmailDomainVerifier.jsx` in your React app.

## ⚠️ **Important Note**

Both versions only verify **DOMAINS** (like gmail.com), not individual email addresses (like user@gmail.com). Individual email verification is impossible due to privacy and security restrictions.