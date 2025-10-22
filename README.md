# Email Verifier

A command-line tool written in Go that verifies email security configurations for domains by checking their DNS records for MX, SPF, and DMARC records.

## ⚠️ **Important: Real User Verification is IMPOSSIBLE**

### **What This Tool CAN Do:**
- ✅ Check if domains exist and can receive emails
- ✅ Verify email security configurations (SPF, DMARC)
- ✅ Validate domain infrastructure

### **What This Tool CANNOT Do:**
- ❌ **Verify if specific users actually exist** (e.g., whether `john@gmail.com` is a real person)
- ❌ **Check if individual email addresses are valid**
- ❌ **Determine if emails will be delivered to specific users**

**Why?** Modern email providers block user verification for privacy and security reasons. No tool can reliably check if individual email addresses like `randomuser@gmail.com` actually exist.

## Overview

This tool helps system administrators, security professionals, and developers verify that domains are properly configured for email security. It performs DNS lookups to check for the presence and content of critical email security records.

## What It Checks

### 1. MX Records (Mail Exchange)
- **Purpose**: Specifies which mail servers can receive email for the domain
- **Importance**: Without MX records, the domain cannot receive emails
- **Output**: Boolean indicating if MX records exist

### 2. SPF Records (Sender Policy Framework)
- **Purpose**: Prevents email spoofing by specifying which servers are authorized to send email for the domain
- **Format**: TXT record starting with "v=spf1"
- **Output**: Boolean indicating presence + full SPF record content

### 3. DMARC Records (Domain-based Message Authentication, Reporting, and Conformance)
- **Purpose**: Provides policy instructions for handling emails that fail SPF/DKIM authentication
- **Format**: TXT record at `_dmarc.domain.com` starting with "v=DMARC1"
- **Output**: Boolean indicating presence + full DMARC record content

#### What It Checks

##### 1. Syntax Validation
- **Purpose**: Ensures email format follows RFC standards
- **Method**: Regular expression pattern matching
- **Example**: `user@domain.com` format validation

##### 2. Domain Existence
- **Purpose**: Verifies the domain part of the email actually exists
- **Method**: DNS host lookup
## Features

- **Batch Processing**: Process multiple domains by entering them one per line
- **CSV Output**: Results are formatted as CSV for easy analysis
- **Error Handling**: Graceful handling of DNS lookup failures
- **Real-time Processing**: Results are displayed immediately as each domain is processed

## Installation

### Prerequisites
- Go 1.15 or higher

### Build from Source
```bash
git clone <repository-url>
cd Email_Verifier
go build -o email-verifier main.go
```

## Usage

### Interactive Mode
Run the program and enter domain names one per line:

```bash
./email-verifier
```

Then type domains:
```
google.com
github.com
example.com
```

Press `Ctrl+C` or `Ctrl+D` to exit.

### Batch Mode with File Input
Process domains from a file:

```bash
./email-verifier < domains.txt
```

### Pipe Input
```bash
echo "google.com" | ./email-verifier
```

## Output Format

The tool outputs results in CSV format with the following columns:

```
domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord
```

### Example Output
```
google.com,true,true,v=spf1 include:_spf.google.com ~all,true,v=DMARC1; p=reject; ruf=mailto:mailauth-reports@google.com
github.com,true,true,v=spf1 ip4:192.30.252.0/22 include:_spf.google.com ~all,true,v=DMARC1; p=reject; rua=mailto:dmarc@github.com
badexample.com,false,false,,false,
```

### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `domain` | string | The domain name being checked |
| `hasMX` | boolean | Whether the domain has MX records |
| `hasSPF` | boolean | Whether the domain has SPF records |
| `spfRecord` | string | Full content of the SPF record (empty if none) |
| `hasDMARC` | boolean | Whether the domain has DMARC records |
| `dmarcRecord` | string | Full content of the DMARC record (empty if none) |

## Use Cases

### Domain Security Auditing
- Verify company domains have proper email security configurations
- Check if domains are vulnerable to email spoofing attacks
- Audit third-party domains before establishing business relationships
- Meet security compliance requirements

### Email Deliverability
- Ensure your domain is properly configured for email sending
- Troubleshoot email delivery issues
- Verify email authentication setup

### Compliance Checking
- Meet security compliance requirements
- Verify DMARC policy implementation
- Document email security posture

### Bulk Domain Analysis
- Process large lists of domains for security assessment
- Generate reports for management
- Monitor domain configurations over time
- **Domain Verification**: Ensure domains exist and can receive emails
- **Infrastructure Testing**: Verify email server configurations
- **Suspicious Pattern Detection**: Flag potentially fake usernames

## Technical Details

### DNS Lookups Performed
1. **MX Lookup**: `net.LookupMX(domain)`
2. **TXT Lookup**: `net.LookupTXT(domain)` for SPF records
3. **DMARC TXT Lookup**: `net.LookupTXT("_dmarc." + domain)`

### Error Handling
- DNS resolution failures are logged but don't stop processing
- Invalid input is handled gracefully
- Network timeouts are managed by Go's default DNS resolver

### Performance
- Concurrent DNS lookups for each record type
- Minimal memory footprint
- Suitable for processing thousands of domains

## Example Domains to Test

Try these domains to see different configurations:

```
google.com          # Well-configured with all records
github.com          # Enterprise-grade configuration
facebook.com        # Social media platform configuration
nonexistentdomain123456.com  # Should show all false values
```

## Troubleshooting

### Common Issues

1. **DNS Resolution Fails**
   - Check internet connectivity
   - Verify domain names are spelled correctly
   - Some domains may block DNS queries from certain IPs

2. **No Output**
   - Ensure you're pressing Enter after each domain
   - Check that Go is properly installed
   - Verify the executable has proper permissions

3. **Partial Results**
   - Some domains may only have MX records but no SPF/DMARC
   - This is common and indicates incomplete email security setup

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license here]

## Future Enhancements

- Add DKIM record checking
- Support for JSON output format
- Parallel processing for faster bulk operations
- Web interface for easier usage
- Historical tracking of domain configurations
- Integration with email security databases

## ⚠️ **Final Note: Why Individual Email Verification is Impossible**

**This tool ONLY checks domains** (like `gmail.com`) **not individual users** (like `john@gmail.com`).

**Technical Reality:**
- No tool can verify if `randomuser@gmail.com` is a real person
- Email providers block user verification for privacy and security
- The only way to verify an email address is to send a verification email
- Any tool claiming to verify individual email addresses has very low accuracy

**What you can do instead:**
1. Use this tool to verify domains are properly configured
2. Send verification emails to users (double opt-in)
3. Monitor email engagement and bounces
4. Accept that some fake emails will get through

**Domain verification = Reliable**  
**User verification = Impossible**