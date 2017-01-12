// Copyright 2016 Jun Zhou <zhoujun06@gmail.com>. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	//	"encoding/json"
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/jzhoucliqr/gitkeeper/config"
	"golang.org/x/oauth2"
)

const GitKeeperToken string = "db414b90c315d8c9523b4080f802e9c3c978c428"

type Config struct {
	Token      string
	Repository *github.Repository
}

type GitObject struct {
	IssuesEvent       *github.IssuesEvent
	PullRequestEvent  *github.PullRequestEvent
	IssueCommentEvent *github.IssueCommentEvent
}

func NewClient(token string) *github.Client {
	flag.CommandLine.Parse([]string{})

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}

func Merge(token string, comment *github.IssueCommentEvent) {
	client := NewClient(token)

	owner := *comment.Issue.User.Login
	repo := *comment.Repo.Name
	number := *comment.Issue.Number

	commitMessage := *comment.Issue.Title

	log.Debugf("try to merged : %s/%s\n", owner, repo)

	result, response, err := client.PullRequests.Merge(owner, repo, number, commitMessage, nil)
	if err != nil {
		log.Debug(err)
	}

	log.Debugf("merged : %b\n", *result.Merged)
	log.Debugf("merge message: %s\n", *result.Message)
	log.Debugf("response : %v\n", *response)

}

// when issue pr, is open, not merged, and is mergeable
func Mergeable(token string, comment *github.IssueCommentEvent) bool {
	if comment.Issue.PullRequestLinks == nil {
		log.Debug("issue is not pr")
		return false
	}
	if *comment.Issue.State == "closed" {
		log.Debug("pr closed")
		return false
	}

	client := NewClient(token)

	owner := *comment.Issue.User.Login
	repo := *comment.Repo.Name
	number := *comment.Issue.Number

	// check if mergable
	p, r, err := client.PullRequests.Get(owner, repo, number)
	if err != nil {
		log.Debug("failed to get pr")
		log.Debugf("response: %v\n", *r)
		return false
	}

	if *p.Merged {
		log.Debug("pr already merged")
		return false
	}
	if p.Mergeable == nil || !*p.Mergeable {
		log.Debug("pr not mergeable")
		return false
	}

	return true
}

func CreateWebHook(keeperConf *config.KeeperConfig, conf *config.RepoConfig) (*github.Hook, error) {
	log.Debugf("conf: %v\n", *conf)

	name := "web"
	owner := conf.User.Name
	repo := conf.Repo.Name
	secret := conf.Repo.WebhookSecret

	events := []string{"pull_request", "issue_comment"}
	c := map[string]interface{}{"content_type": "json", "url": "http://architall.com:8888/api", "secret": secret}

	hook := &github.Hook{Name: &name, Events: events, Config: c}
	//hook := &github.Hook{}
	log.Debugf("hook: %v\n", *hook)

	client := NewClient(conf.User.Token)

	log.Debugf("client: %v\n", *client)

	h, r, err := client.Repositories.CreateHook(owner, repo, hook)

	log.Debugf("response : %v, %v, %v\n", h, r, err)

	if err != nil {
		log.Debug("failed to create webhook")
		return nil, err
	}

	r, err = client.Repositories.AddCollaborator(owner, repo, "keeper-git", nil)
	if err != nil {
		log.Debug("failed to add colab")
	}

	// after add collaborator, an invitation will be sent to keeper-git
	// keeper-git need to accept the invitation
	acceptInvitation(keeperConf.GitUserToken)

	return h, err
}

func acceptInvitation(token string) {
	client := NewClient(GitKeeperToken)

	invitations, r, err := client.Users.ListInvitations()
	if err != nil {
		log.Debugf("failed to list invitation: %v", err)
		return
	}

	for _, invites := range invitations {
		log.Debugf("to accept invitation: %d", *invites.ID)
		r, err = client.Users.AcceptInvitation(*invites.ID)
		if err != nil {
			log.Debug("failed to accetp inviation")
		}

		log.Debugf("res: %v,", r)
	}
}

func ListPR(config Config) {
}
