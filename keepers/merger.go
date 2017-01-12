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

package keepers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	gitobj "github.com/jzhoucliqr/gitkeeper/github"
)

type MergerKeeper struct {
	Ch    chan *gitobj.GitObject
	Token string // user token for git user
}

func (merger *MergerKeeper) Name() string {
	return "Merger"
}

func (merger *MergerKeeper) Receive(obj *gitobj.GitObject) {
	// send obj to channel
	merger.Ch <- obj
}

func (merger *MergerKeeper) KeepDoing(token string) {
	log.Debugf("Keeper %s running", merger.Name())
	merger.Token = token
	for {
		obj := <-merger.Ch
		merger.handle(obj)
	}
}

func (merger *MergerKeeper) handle(obj *gitobj.GitObject) {
	log.Debugf("Keeper: %s handle git object: ", merger.Name())
	if obj.IssuesEvent != nil {
		log.Debugf("issue: %d", *obj.IssuesEvent.Issue.Number)
	} else if obj.PullRequestEvent != nil {
		log.Debugf("pr: %d", *obj.PullRequestEvent.Number)
		//handlePullRequest(obj)
	} else if obj.IssueCommentEvent != nil {
		log.Debugf("comment: %d", *obj.IssueCommentEvent.Issue.Number)
		merger.handlePRComment(obj)
	}
}

func (merger *MergerKeeper) handlePullRequest(obj *gitobj.GitObject) {
}

func (merger *MergerKeeper) handlePRComment(obj *gitobj.GitObject) {
	if obj.IssueCommentEvent == nil {
		return
	}

	comment := obj.IssueCommentEvent

	if gitobj.Mergeable(merger.Token, comment) && mergeable(comment) {
		gitobj.Merge(merger.Token, comment)
	} else {
		log.Debug("not mergeable")
	}
}

func mergeable(comment *github.IssueCommentEvent) bool {
	// if comment includes "LGTM" from admins, then is mergeable
	if *comment.Comment.Body == "LGTM" {
		return true
	}
	return false
}
