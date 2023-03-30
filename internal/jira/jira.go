package jira

import (
	"context"
	"log"

	"github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/tidwall/gjson"
)

type Config struct {
	Host  string
	Email string
	Token string
}

func (j *Config) Create() (*v3.Client, error) {
	atlassian, err := v3.New(nil, j.Host)
	if err != nil {
		return nil, err
	}

	atlassian.Auth.SetBasicAuth(j.Email, j.Token)
	atlassian.Auth.SetUserAgent("curl/7.54.0")

	return atlassian, nil
}

func (j *Config) GetIssueLinks(atlassian *v3.Client, id string) (*models.IssueLinkPageScheme, error) {
	issueLinks, response, err := atlassian.Issue.Link.Gets(context.Background(), id)
	if err != nil {
		if response != nil {
			log.Println("Response HTTP Response", response.Bytes.String())
		}
		panic(err)
	}

	return issueLinks, nil
}

func (j *Config) GetMetadata(atlassian *v3.Client, id string) (gjson.Result, error) {
	metadata, response, err := atlassian.Issue.Metadata.Get(context.Background(), id, false, false)
	if err != nil {
		if response != nil {
			log.Println("Response HTTP Response", response.Bytes.String())
		}
		panic(err)
	}

	return metadata, nil
}

func (j *Config) GetIssue(atlassian *v3.Client, id string) (*models.IssueScheme, error) {
	//https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-get
	issue, response, err := atlassian.Issue.Get(context.Background(), id, []string{"*all"}, []string{"changelog", "schema", "transitions"})

	if err != nil {
		if response != nil {
			log.Println("Response HTTP Response", response.Bytes.String())
		}

		return nil, err
	}

	return issue, nil
}

func (j *Config) QueryJql(atlassian *v3.Client, query string) (*models.IssueSearchScheme, error) {

	var (
		//jql = "project = \"IPD\" AND assignee = currentUser() ORDER BY created DESC"
		jql    = query
		fields = []string{"summary", "status", "project"}
		expand = []string{"changelog", "renderedFields", "names", "schema", "transitions", "operations", "editmeta"}
	)

	issues, response, err := atlassian.Issue.Search.Get(context.Background(), jql, fields, expand, 0, 50, "")

	if err != nil {
		if response != nil {
			log.Println("Response HTTP Response", response.Bytes.String())
		}
		return nil, err
	}

	return issues, nil
}

func (j *Config) GetMyself(atlassian *v3.Client) (*models.UserScheme, error) {
	expand := []string{"groups", "applicationRoles"}
	currentUser, response, err := atlassian.MySelf.Details(context.Background(), expand)
	if err != nil {
		if response != nil {
			log.Println("Response HTTP Response", response.Bytes.String())
		}
		panic(err)
	}

	return currentUser, nil
}

func (j *Config) GetProject(atlassian *v3.Client) (*models.ProjectScheme, error) {
	project, response, err := atlassian.Project.Get(context.Background(), "IPD", []string{"issueTypes"})
	if err != nil {
		if response != nil {
			log.Println("Response HTTP Response", response.Bytes.String())
		}
		panic(err)
	}

	return project, nil
}
