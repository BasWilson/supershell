# supershell

A tiny Go CLI to save and manage SSH connection profiles in a local JSON file.

## Features
- Add, update, delete, list, get connection entries
- Store host, port, user, auth method (key or password), key path or password
- Connect via `ssh` using a saved nickname
- Data stored at: macOS: `~/Library/Application Support/supershell/connections.json`, Linux: `~/.config/supershell/connections.json`

## Install
```bash
# Build and put on your PATH
go build -o bin/supershell ./cmd/supershell
sudo mv bin/supershell /usr/local/bin/
```

Or with Go 1.20+ in this folder:
```bash
go build -o supershell ./cmd/supershell
sudo mv supershell /usr/local/bin/
```

## Usage
```bash
# Add a connection
supershell add --name prod --host 203.0.113.5 --user ubuntu --port 22 --auth key --key $HOME/.ssh/id_rsa

# Add with password (not recommended)
supershell add --name db --host 203.0.113.10 --user root --auth password --password secret

# Update fields
supershell update --name prod --port 2222 --key $HOME/.ssh/id_ed25519

# List all
supershell list

# Get one
supershell get --name prod

# Delete
supershell delete --name db

# Connect using saved settings
supershell connect --name prod
```

## Notes
- Passwords are stored in plain text in the JSON file. Prefer key based auth.
- File and directory permissions default to user only (0700 dir, 0600 file).
- The tool does not manage SSH config; it invokes your system `ssh`.
