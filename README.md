# aStarIsBorn

**Automated GitHub Repository Explorer**  
This project finds a random GitHub repository daily based on configurable criteria and keeps a history in `data/history.jsonl`.

The workflow is fully automated via **GitHub Actions**, but it can also be run locally.

## 🧠 Purpose

- Explore GitHub daily to discover interesting repositories.
- Maintain a JSONL history for visualization or analysis.
- Demonstrate a simple automation setup using Go and GitHub Actions.

## 🔧 Installation

Make sure [Go](https://golang.org/dl/) is installed.

```bash
brew install go
```

Clone the repository:

```bash
git clone https://github.com/<user>/aStarIsBorn.git
cd aStarIsBorn
```

## Project Initialization

```bash
# Initialize Go module
go mod init a-star-is-born

# Download dependencies
go mod tidy
```

## Local Execution

The main script is `run.sh`, which does everything:

- Runs the Go program (./cmd/daily)
- Appends the output JSONL to data/history.jsonl
- Commits and pushes changes to the repository

```bash
# Make the script executable
chmod +x run.sh

# Run the script
./run.sh
```

Detailed logs appear in the console and errors are stored in `error.log`.

You can also run directly using Go:

```bash
# Run without building
go run ./cmd/daily

# Or build a binary
go build -o daily ./cmd/daily
./daily
```

## Automated Execution via GitHub Actions

The workflow `.github/workflows/cron.yml` automatically triggers the script every day at 9 AM UTC.

- New repositories found are added to `data/history.jsonl`.
- Errors are stored in `error.log`.
- Commits are automatically pushed to the repository.

No manual intervention is needed to keep the history up-to-date.

## 📂 Repository Structure

- `cmd/daily/main.go` – main code for fetching and processing GitHub repositories
- `run.sh` – complete execution script (local or CI)
- `data/history.jsonl` – daily history of repositories found
- `error.log `– log file for errors
- `.github/workflows/cron.yml` – GitHub Actions workflow

## Roadmap / Ideas for improvements

- [x] 📧 Send daily email with the repository information of the day
- [ ] 🔄 Avoid duplicates: skip repositories already in history.jsonl
- [ ] 🤖 Add AI to make a summary of the repo and rebuild a mini README.md
- [ ] 🌐 Use this as base to build a Contributor Territory Map
- [ ] ...
