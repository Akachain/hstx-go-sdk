# Akachain - High Secure Transaction Samples ⌁ hstx-chaincode

[![Go Report Card](https://goreportcard.com/badge/github.com/Akachain/hstx-go-sdk)](https://goreportcard.com/report/github.com/Akachain/hstx-go-sdk)

Chaincode (Smart Contract)

## Prerequisites
- OS: Ubuntu 18.04, Mac OS 10.13.6
- Language: Go 1.13+
- IDE: Visual Studio Code + Go plugin / IntelliJ GoLand
- Unit test tool: built-in testing command (go test)

## Require

To create SuperAdmin and Approval, the invoking identity must include attribute "hstx.role=SuperAdmin" in it's certificate.

Example create certificate with role SuperAdmin

```
curl --location --request POST 'http://admin-service-address/registerUser' \
--header 'Content-Type: application/json' \
--data-raw '{
	"orgname": 	"operator",
	"username": "SuperAdmin",
	"maxEnrollments": 1,
	"attrs": [{ "name": "hstx.role", "value": "SuperAdmin", "ecert": true }]
}'
```