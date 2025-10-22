package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// DomainCheckResult represents the result of domain verification
type DomainCheckResult struct {
	Domain      string `json:"domain"`
	HasMX       bool   `json:"hasMX"`
	HasSPF      bool   `json:"hasSPF"`
	SPFRecord   string `json:"spfRecord"`
	HasDMARC    bool   `json:"hasDMARC"`
	DMARCRecord string `json:"dmarcRecord"`
	Error       string `json:"error,omitempty"`
	CheckedAt   string `json:"checkedAt"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    string `json:"time"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func main() {
	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	api.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// Domain verification endpoint
	api.HandleFunc("/verify/domain/{domain}", verifyDomainHandler).Methods("GET")

	// Bulk domain verification endpoint
	api.HandleFunc("/verify/domains", verifyDomainsHandler).Methods("POST")

	// Setup CORS to allow frontend connections
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // In production, specify your frontend domain
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Wrap router with CORS
	handler := c.Handler(r)

	// Start server
	port := "8080"
	fmt.Printf("🚀 Email Verifier API starting on port %s\n", port)
	fmt.Printf("📋 API Documentation:\n")
	fmt.Printf("   GET  /api/v1/health                    - Health check\n")
	fmt.Printf("   GET  /api/v1/verify/domain/{domain}    - Verify single domain\n")
	fmt.Printf("   POST /api/v1/verify/domains            - Verify multiple domains\n")
	fmt.Printf("\n🌐 Example usage:\n")
	fmt.Printf("   curl http://localhost:%s/api/v1/verify/domain/gmail.com\n", port)
	fmt.Printf("\n⚠️  Note: This API only verifies DOMAINS, not individual email addresses!\n")

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// verifyDomainHandler handles single domain verification
func verifyDomainHandler(w http.ResponseWriter, r *http.Request) {
	// Extract domain from URL path
	vars := mux.Vars(r)
	domain := vars["domain"]

	// Validate domain parameter
	if domain == "" {
		sendErrorResponse(w, "Domain parameter is required", http.StatusBadRequest)
		return
	}

	// Clean and validate domain
	domain = strings.TrimSpace(strings.ToLower(domain))
	if !isValidDomain(domain) {
		sendErrorResponse(w, "Invalid domain format", http.StatusBadRequest)
		return
	}

	// Perform domain check
	result := checkDomainAPI(domain)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// verifyDomainsHandler handles bulk domain verification
func verifyDomainsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestBody struct {
		Domains []string `json:"domains"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		sendErrorResponse(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate domains list
	if len(requestBody.Domains) == 0 {
		sendErrorResponse(w, "At least one domain is required", http.StatusBadRequest)
		return
	}

	if len(requestBody.Domains) > 50 {
		sendErrorResponse(w, "Maximum 50 domains allowed per request", http.StatusBadRequest)
		return
	}

	// Process each domain
	var results []DomainCheckResult
	for _, domain := range requestBody.Domains {
		domain = strings.TrimSpace(strings.ToLower(domain))
		if isValidDomain(domain) {
			result := checkDomainAPI(domain)
			results = append(results, result)
		} else {
			// Add error result for invalid domain
			result := DomainCheckResult{
				Domain:    domain,
				Error:     "Invalid domain format",
				CheckedAt: time.Now().UTC().Format(time.RFC3339),
			}
			results = append(results, result)
		}
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
		"total":   len(results),
	})
}

// checkDomainAPI performs domain verification and returns structured result
func checkDomainAPI(domain string) DomainCheckResult {
	result := DomainCheckResult{
		Domain:    domain,
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// 1. CHECK MX RECORDS
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		result.Error = fmt.Sprintf("MX lookup failed: %v", err)
	} else {
		result.HasMX = len(mxRecords) > 0
	}

	// 2. CHECK SPF RECORDS
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		if result.Error == "" {
			result.Error = fmt.Sprintf("TXT lookup failed: %v", err)
		}
	} else {
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=spf") {
				result.HasSPF = true
				result.SPFRecord = txt
				break
			}
		}
	}

	// 3. CHECK DMARC RECORDS
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		// DMARC lookup failure is not critical, don't overwrite existing errors
		log.Printf("DMARC lookup failed for %s: %v", domain, err)
	} else {
		for _, txt := range dmarcRecords {
			if strings.HasPrefix(txt, "v=DMARC") {
				result.HasDMARC = true
				result.DMARCRecord = txt
				break
			}
		}
	}

	return result
}

// isValidDomain performs basic domain format validation
func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// Basic check for domain format
	if strings.Contains(domain, " ") || strings.Contains(domain, "\t") {
		return false
	}

	// Must contain at least one dot
	if !strings.Contains(domain, ".") {
		return false
	}

	// Should not start or end with dot or hyphen
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") ||
		strings.HasPrefix(domain, "-") || strings.HasSuffix(domain, "-") {
		return false
	}

	return true
}

// sendErrorResponse sends a JSON error response
func sendErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResp := ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Code:    code,
	}

	json.NewEncoder(w).Encode(errorResp)
}
