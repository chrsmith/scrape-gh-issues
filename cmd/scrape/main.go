package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	repoSlug := flag.String("repo", "", "Name of the GitHub repo, e.g. chrsmith/scrape-gh-issues.")

	flag.Parse()

	if repoSlug == nil {
		fmt.Println("Required flag repo not set.")
		os.Exit(1)
	}
	parts := strings.Split(*repoSlug, "/")
	if len(parts) != 2 {
		fmt.Println("Invalid repo slug.")
		os.Exit(1)
	}
	owner := parts[0]
	repo := parts[1]

	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if accessToken == "" {
		fmt.Println("Please specify a GitHub access token in the GITHUB_ACCESS_TOKEN env variable.")
		os.Exit(1)
	}

	// Create GitHub client.
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	client.UserAgent = "scrape-gh-issues"

	// Scrape
	var allIssues []*github.Issue
	page := 0
	for {
		opts := &github.IssueListByRepoOptions{
			State: "open",
		}
		opts.Page = page
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		allIssues = append(allIssues, issues...)

		page = resp.NextPage
		if page == 0 {
			break
		}
	}

	fmt.Println("Number, Title, Milestone, Labels")
	for _, issue := range allIssues {
		// Put titles in quotes, escaping double quote characters.
		title := fmt.Sprintf("%q", strings.Replace(issue.GetTitle(), "\"", "\"\"", -1))

		milestone := ""
		if issue.Milestone != nil {
			// Avoid "0.10" being interpreted as float 0.1.
			milestone = "M" + issue.Milestone.GetTitle()
		}
		labels := ""
		for _, label := range issue.Labels {
			labels = label.GetName() + ";" + labels
		}
		fmt.Printf("%d,%s,%s,%v\n", issue.GetNumber(), title, milestone, labels)
	}

	os.Exit(0)
}
