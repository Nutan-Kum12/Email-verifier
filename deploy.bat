@echo off
echo ==========================================
echo  EMAIL VERIFIER - BUILD AND DEPLOY
echo ==========================================
echo.

echo What would you like to build?
echo 1. API Server (for web/frontend integration)
echo 2. CLI Tool (for command-line usage)
echo 3. Both
echo.
set /p choice="Enter your choice (1, 2, or 3): "

if "%choice%"=="1" goto build_api
if "%choice%"=="2" goto build_cli
if "%choice%"=="3" goto build_both
echo Invalid choice. Building API server by default...
goto build_api

:build_api
echo.
echo ======================================
echo Building API Server...
echo ======================================
echo.

echo Step 1: Installing Go dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo ERROR: Failed to install dependencies
    echo Make sure Go is installed and you have internet connection
    pause
    exit /b 1
)
echo ✅ Dependencies installed successfully

echo.
echo Step 2: Building API server...
go build -o email-verifier-api.exe api_server.go
if %errorlevel% neq 0 (
    echo ERROR: Failed to build API server
    pause
    exit /b 1
)
echo ✅ API server built successfully

echo.
echo Step 3: Starting API server...
start /B email-verifier-api.exe
echo ⏳ Waiting for server to start...
timeout /t 3 /nobreak >nul

echo.
echo ==========================================
echo 🚀 API SERVER READY!
echo ==========================================
echo.
echo Your Email Verifier API is running on:
echo 👉 http://localhost:8080
echo.
echo API Endpoints:
echo • GET  /api/v1/health
echo • GET  /api/v1/verify/domain/{domain}
echo • POST /api/v1/verify/domains
echo.
echo Frontend Example:
echo 👉 Open frontend_example.html in your browser
echo.
echo Test the API:
echo 👉 curl http://localhost:8080/api/v1/verify/domain/gmail.com
echo.
goto end

:build_cli
echo.
echo ======================================
echo Building CLI Tool...
echo ======================================
echo.

echo Step 1: Building CLI tool...
cd cli
go build -o email-verifier-cli.exe main.go
if %errorlevel% neq 0 (
    echo ERROR: Failed to build CLI tool
    pause
    exit /b 1
)
echo ✅ CLI tool built successfully
move email-verifier-cli.exe ..
cd ..

echo.
echo ==========================================
echo 🚀 CLI TOOL READY!
echo ==========================================
echo.
echo Your Email Verifier CLI tool is ready!
echo.
echo Usage Examples:
echo • echo gmail.com ^| email-verifier-cli.exe
echo • email-verifier-cli.exe ^< domains.txt
echo.
echo Test the CLI:
echo 👉 echo gmail.com ^| email-verifier-cli.exe
echo.
goto end

:build_both
echo.
echo ======================================
echo Building Both API Server and CLI Tool...
echo ======================================
echo.

echo Step 1: Installing Go dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo ERROR: Failed to install dependencies
    pause
    exit /b 1
)
echo ✅ Dependencies installed

echo.
echo Step 2: Building API server...
go build -o email-verifier-api.exe api_server.go
if %errorlevel% neq 0 (
    echo ERROR: Failed to build API server
    pause
    exit /b 1
)
echo ✅ API server built

echo.
echo Step 3: Building CLI tool...
cd cli
go build -o email-verifier-cli.exe main.go
if %errorlevel% neq 0 (
    echo ERROR: Failed to build CLI tool
    pause
    exit /b 1
)
move email-verifier-cli.exe ..
cd ..
echo ✅ CLI tool built

echo.
echo Step 4: Starting API server...
start /B email-verifier-api.exe
echo ⏳ Waiting for server to start...
timeout /t 3 /nobreak >nul

echo.
echo ==========================================
echo 🚀 BOTH TOOLS READY!
echo ==========================================
echo.
echo API Server:
echo 👉 http://localhost:8080
echo 👉 Open frontend_example.html in browser
echo.
echo CLI Tool:
echo 👉 echo gmail.com ^| email-verifier-cli.exe
echo.

:end
echo To stop the API server, press Ctrl+C in the terminal
echo.
pause