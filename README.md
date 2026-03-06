# Gender-API Command Line Client

This directory contains the official Command Line Client for [Gender-API.com](https://gender-api.com). It is built in Go and allows you to interact with all major API v2 endpoints directly from your preferred terminal environment.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or higher recommended)
- A valid Gender-API account and API token.

## Compiling the Client

To build the client natively for your current operating system, simply run:

```bash
go build -o gender-api-cli main.go
```

This will produce a standalone executable binary named `gender-api-cli` inside this directory.

### Cross-Compilation (Linux, macOS, Windows)

We have provided a `build.sh` script that automatically cross-compiles the CLI for all major platforms (amd64 and arm64 arrays). The resulting executables will be placed in the `build/` directory.

```bash
chmod +x build.sh
./build.sh
```

#### Compile using Docker
If you do not have Go installed locally, you can use the official `golang` Docker container to run the build script securely. We have provided a wrapper script in the root `bin` directory of the Gender-API project to make this effortless.

From the project root, simply run:
```bash
./bin/build-cli-docker.sh
```

## Configuration

Before making any queries, you must provide your Gender-API authorization token. The CLI client checks for your API key in the following order of precedence:

1. **Command Line Flag:** Pass the key directly using the `-key` flag.
   ```bash
   ./gender-api-cli -key "your-api-key-here" -first_name "Sandra"
   ```

2. **Environment Variable:** Set the `GENDER_API_KEY` environment variable. This is highly recommended to prevent your key from being logged in shell history.
   ```bash
   export GENDER_API_KEY="your-api-key-here"
   ```

3. **Config File:** Create a file named `.gender-api-key` in your user's home directory (`~/.gender-api-key`) and paste your API key inside it (plain text, no quotes).
   ```bash
   echo "your-api-key-here" > ~/.gender-api-key
   ```

## macOS Installation Note (Gatekeeper)

If you download the pre-compiled binary on macOS using your browser, macOS's Gatekeeper feature will likely flag it with `com.apple.quarantine` causing an "Apple could not verify it is free of malware" error.

Because our CLI is an open-source binary, you can whitelist it quickly by opening your terminal and removing the quarantine extended attribute:
```bash
xattr -d com.apple.quarantine /path/to/gender-api-cli-darwin-arm64
```

## Usage and Examples

You can view the full list of available flags at any time by running:
```bash
./gender-api-cli -h
```

### 1. Standard Name Queries

#### Query by First Name
```bash
./gender-api-cli -first_name "Sandra" -country "US"
```

#### Query by Full Name
```bash
./gender-api-cli -full_name "Theresa Miller"
```

#### Query by Email Address
```bash
./gender-api-cli -email "thomasfreeman@example.com"
```

### 2. Country of Origin Queries
To determine likely origins instead of just the gender, append the `-origin` flag.

```bash
./gender-api-cli -first_name "Sandra" -origin
```

### 3. Account Statistics
To check your remaining credits without consuming a new credit:

```bash
./gender-api-cli -stats
```

## Pipeline Integration (JSON Output)

By default, the client formats the response in an easy-to-read text block. If you are building automated scripts or feeding data into another application (like `jq` or an AI tool), use the `-out=json` flag.

```bash
./gender-api-cli -first_name "Alex" -out=json | jq .
```

## Running the Automated Test Script

We have provided a small `test.sh` bash script in this directory that demonstrates compiling the program and running a batch of automated sanity tests against the actual API.

**Note:** Executing this script will consume real API credits.

```bash
# Make it executable if it isn't already
chmod +x test.sh

# Run the tests
./test.sh
```

## Developer: Mirroring to GitHub

The source code for this CLI client lives inside the main Gender-API project repository (`/cli-client/`). However, we independently mirror this directory to a public GitHub repository at [markus-perl/gender-api-cli](https://github.com/markus-perl/gender-api-cli) for users to review the open-source code.

To publish your latest changes from the main repository to the public GitHub mirror, use the `git subtree` script provided in the root `bin/` directory:

```bash
cd /path/to/gender-api
./bin/push-cli-to-github.sh
```

This script will safely extract the `cli-client/` directory commits and push them to the `main` branch of the external repository without including the rest of the Gender-API server codebase or the compiled `build/` files (which are git-ignored).
