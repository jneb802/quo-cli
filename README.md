# quo-cli

A single-binary CLI for the OpenPhone messaging API. Send and list SMS messages from the terminal.

## Install

```bash
go install github.com/jneb802/quo-cli@latest
```

Or build from source:

```bash
git clone https://github.com/jneb802/quo-cli.git
cd quo-cli
go build -o quo .
```

## Setup

Set your API key:

```bash
export QUO_API_KEY="your-openphone-api-key"
```

Optional defaults to avoid repeating flags:

```bash
export QUO_FROM="+15551234567"           # default --from number
export QUO_PHONE_NUMBER_ID="PNxxxxxx"    # default --phone-number-id
```

## Usage

### Send a message

```bash
quo send --from "+15551234567" --to "+15559876543" "Hello world"

# with QUO_FROM set:
quo send --to "+15559876543" "Hello world"
```

### List messages

```bash
quo list --phone-number-id "PNxxx" --participant "+15559876543"
quo list --phone-number-id "PNxxx" --participant "+15559876543" --limit 50
```

### Get a message

```bash
quo get AC123abc
```

### JSON output

Add `--json` to any command for machine-readable output:

```bash
quo send --to "+15559876543" "Hello" --json
quo list --phone-number-id "PNxxx" --participant "+15559876543" --json
```
