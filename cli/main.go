// CLI Version of Email Verifier
// This is the original command-line interface version
//
// To build and use this CLI version:
// cd cli
// go build -o email-verifier-cli.exe main.go
//
// Usage:
// echo "gmail.com" | email-verifier-cli.exe
// email-verifier-cli.exe < domains.txt
//
// The CLI version reads domains from stdin and outputs CSV results
// The API version (../api_server.go) provides REST endpoints for web integration

package main

import (
	"bufio"   // Package for buffered I/O operations, used to read input line by line
	"fmt"     // Package for formatted I/O operations (printing to console)
	"log"     // Package for logging errors and messages
	"net"     // Package for network operations, used for DNS lookups
	"os"      // Package for operating system interface, used to read from stdin
	"strings" // Package for string manipulation operations
)

func main() {
	// Create a scanner to read input from standard input (stdin) line by line
	scanner := bufio.NewScanner(os.Stdin)

	// Print CSV header row to define the output format
	fmt.Printf("domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord\n")

	// Loop through each line of input until EOF or error
	for scanner.Scan() {
		// Process each domain name entered by the user
		checkDomain(scanner.Text())
	}

	// Check if there was an error while reading from stdin
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading standard input: %v", err)
	}
}

func checkDomain(domain string) {
	// Initialize variables to store the results of our DNS checks
	var hasMX, hasSPF, hasDMARC bool  // Boolean flags for each record type
	var spfRecord, dmarcRecord string // String variables to store the actual record content

	// 1. CHECK MX RECORDS
	// MX (Mail Exchange) records specify which mail servers can receive email for this domain
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		// Log error if MX record lookup fails (domain might not exist or DNS issues)
		log.Printf("Error looking up MX records for %s: %v", domain, err)
	} else {
		// Domain has MX records if any records were returned
		hasMX = len(mxRecords) > 0
	}

	// 2. CHECK SPF RECORDS
	// SPF (Sender Policy Framework) records are stored in TXT records and help prevent email spoofing
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		// Log error if TXT record lookup fails
		log.Printf("Error looking up TXT records for %s: %v", domain, err)
	} else {
		// Loop through all TXT records to find SPF record
		for _, txt := range txtRecords {
			// SPF records start with "v=spf1" (version 1 is the current standard)
			// Commented out old approach that checked first 4 characters manually
			// if len(txt) >= 4 && txt[:4] == "v=spf" {
			// 	hasSPF = true
			// 	spfRecord = txt
			// }

			// More robust approach using strings.HasPrefix
			if strings.HasPrefix(txt, "v=spf") {
				hasSPF = true
				spfRecord = txt // Store the entire SPF record
				break           // Exit loop once we find the SPF record
			}
		}
	}

	// 3. CHECK DMARC RECORDS
	// DMARC (Domain-based Message Authentication, Reporting, and Conformance) records
	// are stored in TXT records at the subdomain "_dmarc"
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		// Log error if DMARC record lookup fails
		log.Printf("Error looking up DMARC records for %s: %v", domain, err)
	} else {
		// Loop through TXT records at _dmarc subdomain to find DMARC policy
		for _, txt := range dmarcRecords {
			// DMARC records start with "v=DMARC1"
			if strings.HasPrefix(txt, "v=DMARC") {
				hasDMARC = true
				dmarcRecord = txt // Store the entire DMARC record
				break             // Exit loop once we find the DMARC record
			}
		}
	}

	// OUTPUT RESULTS
	// Print results in CSV format: domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord
	fmt.Printf("%s,%t,%t,%s,%t,%s\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
