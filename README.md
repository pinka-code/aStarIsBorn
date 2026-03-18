# aStarIsBorn

```bash
# install
brew install go

# init
go mod init repo-selector
go mod tidy

# run
cd repo-selector
go run ./cmd/daily

# build
cd repo-selector
go build -o daily ./cmd/daily
./daily
```