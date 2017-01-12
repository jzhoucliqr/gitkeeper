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

// storage package provide persistence fro gitkeeper
// will support both etcd and consul

package kvstore

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"github.com/google/go-github/github"
	"github.com/jzhoucliqr/gitkeeper/config"
	gitobj "github.com/jzhoucliqr/gitkeeper/github"
	"github.com/jzhoucliqr/gitkeeper/keepers"
	"strings"
	"time"
)

func NewKVStore(conf *config.KeeperConfig) (store.Store, error) {
	etcd.Register()
	backend := store.ETCD
	backendUrl := conf.KVStoreUrl

	kvstore, err := libkv.NewStore(
		backend,
		[]string{backendUrl},
		&store.Config{ConnectionTimeout: 10 * time.Second},
	)

	return kvstore, err
}

func ConsumeEvents(conf *config.KeeperConfig) {
	kvstore, err := NewKVStore(conf)
	if err != nil {
		log.Fatalf("%v", err)
	}

	opts := &store.WatchOptions{Recursive: true, ReturnAll: false}
	stopCh := make(<-chan struct{})
	events, err := kvstore.WatchChild("/", stopCh, opts)

	kps := keepers.GetEnabledKeepers()
	for {
		select {
		case <-stopCh:
			return
		case pairs := <-events:
			for _, pair := range pairs {
				log.Debugf("value changed on key %v: ", pair.Key)
				obj := &gitobj.GitObject{}

				if strings.HasPrefix(pair.Key, "/prs") {
					pr := &github.PullRequestEvent{}
					if err := json.Unmarshal(pair.Value, pr); err != nil {
						log.Error(err)
					}
					obj.PullRequestEvent = pr
				} else if strings.HasPrefix(pair.Key, "/comments") {
					comment := &github.IssueCommentEvent{}
					if err := json.Unmarshal(pair.Value, comment); err != nil {
						log.Error(err)
					}
					obj.IssueCommentEvent = comment
				} else if strings.HasPrefix(pair.Key, "/issues") {
					issue := &github.IssuesEvent{}
					if err := json.Unmarshal(pair.Value, issue); err != nil {
						log.Error(err)
					}
					obj.IssuesEvent = issue
				} else {
					log.Error("invalid key: " + pair.Key)
				}

				for _, keeper := range kps {
					log.Debug("send obj to keeper")
					keeper.Receive(obj)
				}
			}
		}
	}
}
