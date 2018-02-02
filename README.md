# scrape-gh-issues

Export a GitHub repos issue's to a CSV file.

Requires a GitHub personal access token. You can create one here:
https://github.com/settings/tokens

## Installation

```
go get -u github.com/chrsmith/scrape-gh-issues/cmd/scrape

# Scrape issues from chrsmith/scrape-gh-issues.
export GITHUB_ACCESS_TOKEN="..."
scrape --repo chrsmith/scrape-gh-issue
```
