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
	"github.com/jzhoucliqr/gitkeeper/github"
)

// ClaKeeper enforce cla:yes
type ClaKeeper struct {
	Ch chan *github.GitObject
}

func (cla *ClaKeeper) Name() string {
	return "CLA"
}

func (cla *ClaKeeper) Receive(obj *github.GitObject) {
	// send obj to channel
	cla.Ch <- obj
}

func (cla *ClaKeeper) KeepDoing(token string) {
	log.Debugf("Keeper %s running", cla.Name())
	for {
		obj := <-cla.Ch
		cla.handle(obj)
	}
}

func (cla *ClaKeeper) handle(obj *github.GitObject) {
	log.Debugf("Keeper: %s handle git object: ", cla.Name())
	if obj.IssuesEvent != nil {
		log.Debugf("issue: %d", *obj.IssuesEvent.Issue.Number)
	} else if obj.PullRequestEvent != nil {
		log.Debugf("pr: %d", *obj.PullRequestEvent.Number)
	} else if obj.IssueCommentEvent != nil {
		log.Debugf("comment: %d", *obj.IssueCommentEvent.Issue.Number)
	}
}
