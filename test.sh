#!/bin/bash

# A simple script to compile the gender-api-cli and run basic sanity checks.
# Note: This will consume a few API credits from the account associated with the GENDER_API_KEY.

set -e

echo "==> Compiling CLI Tool..."
go build -o gender-api-cli main.go
echo "==> Compilation successful."
echo ""

# Check if auth works by running a simple stats query
if ! ./gender-api-cli -stats > /dev/null 2>&1; then
    echo "ERROR: Could not authenticate with Gender-API."
    echo "Please provide your API key via:"
    echo "  1. GENDER_API_KEY environment variable"
    echo "  2. ~/.gender-api-key file"
    echo "  3. Modify this script to pass -key"
    exit 1
fi

echo "==> Test 1: Query by First Name (Text Output)"
./gender-api-cli -first_name "Sandra" -country "US"
echo "-----------------------------------"
echo ""

echo "==> Test 2: Query by Full Name (JSON Output)"
./gender-api-cli -full_name "Thomas Johnson" -out=json
echo "-----------------------------------"
echo ""

echo "==> Test 3: Query Country of Origin"
./gender-api-cli -first_name "Johann" -origin
echo "-----------------------------------"
echo ""

echo "==> Test 4: Query Account Statistics"
./gender-api-cli -stats
echo "-----------------------------------"
echo ""

echo "==> All automated tests completed."
